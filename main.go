package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/redkite1/zigbee-gw/mqtt"
	"github.com/redkite1/zigbee-gw/xbee"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")

	//TODO What about giving an argument on startup for specifying config-dir?
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath("./etc")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	log.Infoln("Using config:", viper.ConfigFileUsed())

	viper.UnmarshalKey("zb_sources", &registeredZBSources)
	for k, v := range registeredZBSources {
		log.Infoln("New ZigBee source:", k, v)
	}

	if viper.GetBool("debug_mode") {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	// Prepare ZigBee listenning
	ZBframeChan := make(chan xbee.APIFrame)
	ZBstopped := make(chan bool)

	xbee.InitSerial(viper.GetString("serial.name"), viper.GetInt("serial.speed"))
	mqtt.InitMQTT(viper.GetString("mqtt.host"), viper.GetInt("mqtt.port"), viper.GetString("mqtt.username"), viper.GetString("mqtt.password"))

	go xbee.ReadSerial(ZBframeChan)
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

	// Disconnect from MQTT
	mqtt.Client.Disconnect(5000)
	// Stop processing more ZigBee frames
	close(ZBframeChan)
	// Stop processing more TCP requests
	TCPstop <- true

	<-ZBstopped
	<-TCPstopped

	log.Infoln("Application stopped gracefully")
}
