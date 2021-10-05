package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/redkite1/zigbee-gw/xbee"

	"github.com/oklog/run"
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
	var g run.Group

	// receiving and process ZigBee frames
	{
		ZBframeChan := make(chan xbee.APIFrame)
		xbee.InitSerial(viper.GetString("serial.name"), viper.GetInt("serial.speed"))
		go xbee.ReadSerial(ZBframeChan)

		//mqtt.InitMQTT(viper.GetString("mqtt.host"), viper.GetInt("mqtt.port"), viper.GetString("mqtt.username"), viper.GetString("mqtt.password"))

		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			processZBFrames(ctx, ZBframeChan)
			return nil
		}, func(error) {
			cancel()
		})
	}

	// receive and process HTTP requests
	{
		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			processTCPrequests(ctx)
			return nil
		}, func(error) {
			cancel()
		})
	}

	// handle termination on signals
	{
		cancelHandlerChan := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelHandlerChan:
				log.Infof("interrupted... no more signals to handle")
				return nil
			}
		}, func(error) {
			close(cancelHandlerChan)
		})
	}

	log.Infof("application stopped gracefully: %v", g.Run())
}
