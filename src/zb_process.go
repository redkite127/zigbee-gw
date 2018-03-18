package main

import (
	"encoding/hex"
	"log"

	"github.com/redkite1/zigbee-gw/src/xbee"
)

func processZBFrames(fc <-chan xbee.Frame, stop chan<- bool) {
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
	stop <- true
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
