package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/redkite1/zigbee-gw/xbee"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func processZBFrames(fc <-chan xbee.ReceivePacketFrame, stopped chan<- bool) {
	var err error
	log.Infof("waiting ZigBee frames: %s", viper.GetString("serial.name"))
	for f := range fc {
		switch f.Type {
		case xbee.TypeReceivePacket:
			err = processReceivePacketFrame(f)
		case xbee.TypeRemoteATCommandResponse:
			err = processRemoteATCommandResponseFrame(f)
		case xbee.TypeTransmitStatus:
			continue
		default:
			log.Printf("Unsupported frame type: %X\n", f.Type)
		}
		if err != nil {
			log.Errorln(err)
		}
		log.Debugf("==============================================================")
	}
	log.Infof("interrupted... no more ZigBee frame to process")
	stopped <- true
}

func processReceivePacketFrame(f xbee.ReceivePacketFrame) error {
	sa64 := f.Data[:8]
	log.Debugf("64-bit source address: % X", sa64)

	sa16 := f.Data[8:10]
	log.Debugf("16-bit source address: % X", sa16)

	ro := f.Data[10:11]
	log.Debugf("Receive options: % X", ro)

	rfd := f.Data[11:]
	log.Debugf("RF data: % X", rfd)
	log.Debugf("RF data (string): %s", rfd)

	if err := ZBredirect(hex.EncodeToString(sa64), hex.EncodeToString(sa16), rfd); err != nil {
		return err
	}

	return nil
}

func processRemoteATCommandResponseFrame(f xbee.ReceivePacketFrame) error {
	var frame xbee.RemoteATCommandResponseFrame

	//TODO frame.StartDelimiter
	//TODO frame.Length
	frame.Type = byte(f.Type)
	frame.ID = f.Data[0]
	copy(frame.SourceAddress64[:], f.Data[1:9])
	copy(frame.SourceAddress16[:], f.Data[9:11])
	copy(frame.ATCommand[:], f.Data[11:13])
	frame.CommandStatus = f.Data[13]

	frame.ParameterValue = make([]byte, len(f.Data[14:]))
	copy(frame.ParameterValue[:], f.Data[14:])

	//TODO frame.Checksum

	switch strings.ToUpper(string(frame.ATCommand[:])) {
	case "SH":
		xbee.RecordSH(fmt.Sprintf("%X", frame.SourceAddress16), fmt.Sprintf("%X", frame.ParameterValue))
	case "SL":
		xbee.RecordSL(fmt.Sprintf("%X", frame.SourceAddress16), fmt.Sprintf("%X", frame.ParameterValue))
	}

	return nil
}
