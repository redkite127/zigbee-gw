package main

import (
	"encoding/hex"

	"github.com/redkite1/zigbee-gw/src/xbee"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func processZBFrames(fc <-chan xbee.Frame, stopped chan<- bool) {
	var err error
	log.Printf("Waiting ZigBee frames: %s", viper.GetString("serial.name"))
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
	log.Printf("No more ZigBee frame to process")
	stopped <- true
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

	if err := ZBredirect(hex.EncodeToString(sa64), hex.EncodeToString(sa16), rfd); err != nil {
		return err
	}

	return nil
}
