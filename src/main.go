package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/redkite1/zigbee-gw/src/xbee"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")

	//TODO What about giving an argument on startup for specifying config-dir?
	viper.AddConfigPath(os.Getenv("etc_dir")) //TODO Handle ENV variable with viper
	viper.AddConfigPath("/opt/zigbee-gw/etc")
	viper.AddConfigPath("./etc")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../etc")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	log.Infoln("Using config:", viper.ConfigFileUsed())

	viper.UnmarshalKey("zb_sources", &registeredZBSources)
	for k, v := range registeredZBSources {
		log.Infoln("New ZigBee source:", k, v)
	}
}

func main() {
	// Prepare ZigBee listenning
	ZBframeChan := make(chan xbee.Frame)
	ZBstopped := make(chan bool)

	go xbee.ReadSerial(ZBframeChan, viper.GetString("serial.name"), viper.GetInt("serial.speed"))
	go processZBFrames(ZBframeChan, ZBstopped)

	// Prepare TCP listenning
	TCPstop := make(chan bool)
	TCPstopped := make(chan bool)

	go processTCPrequests(TCPstop, TCPstopped)

	// Prepare graceful shutdown
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	sig := <-gracefulStop
	log.Infof("Caught signal: %+v", sig)
	log.Info("Stopping gracefully the application...")

	// Stop processing more ZigBee frames
	close(ZBframeChan)
	// Stop processing more TCP requests
	TCPstop <- true

	<-ZBstopped
	<-TCPstopped

	log.Infoln("Application stopped gracefully")
}
