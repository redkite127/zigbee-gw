package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/redkite1/zigbee-gw/src/xbee"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func processFrames(fc <-chan xbee.Frame) {
	var err error
	for f := range fc {
		switch f.Type {
		case xbee.TypeReceivePacket:
			err = processReceivePacketFrame(f)
		default:
			log.Printf("Unsupported frame type: %X\n", f.Type)
		}
		if err != nil {
			log.Println(err)
		}
		log.Println("==============================================================")
	}
}

func processReceivePacketFrame(f xbee.Frame) error {
	sa64 := f.Data[:8]
	log.Printf("64-bit source address: % X", sa64)

	sa16 := f.Data[8:10]
	log.Printf("16-bit source address: % X", sa16)

	ro := f.Data[10:11]
	log.Printf("Receive options: % X", ro)

	rfd := f.Data[11:]
	log.Printf("RF data: % X", rfd)
	log.Printf("RF data (string): %s", rfd)

	if err := redirect(hex.EncodeToString(sa64), rfd); err != nil {
		return err
	}

	return nil
}

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
	//stopChan := make(chan interface{})
	frameChan := make(chan xbee.Frame)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		close(frameChan)
	}()

	go xbee.ReadSerial(frameChan, viper.GetString("serial.name"), viper.GetInt("serial.speed"))
	processFrames(frameChan)
}
