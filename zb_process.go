package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/redkite1/zigbee-gw/mqtt"
	"github.com/redkite1/zigbee-gw/xbee"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func processZBFrames(fc <-chan xbee.APIFrame, stopped chan<- bool) {
	var err error
	log.Infof("waiting ZigBee frames: %s", viper.GetString("serial.name"))
	for f := range fc {
		switch xbee.ReceiveAPIFrameType(f.Data[0]) {
		case xbee.TypeReceivePacket:
			err = processReceivePacketFrame(f)
		case xbee.TypeRemoteATCommandResponse:
			err = processRemoteATCommandResponseFrame(f)
		case xbee.TypeTransmitStatus:
			continue
		default:
			log.Printf("Unsupported frame type: %X\n", f.Data[0])
		}
		log.Debugln(f)
		if err != nil {
			log.Errorln(err)
		}
		log.Debugf("==============================================================")
	}
	log.Infof("interrupted... no more ZigBee frame to process")
	stopped <- true
}

func processReceivePacketFrame(f xbee.APIFrame) error {
	//TODO move that in xbee/api_frame_data_receive_packet.go in a FromBytes functions?

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

func processRemoteATCommandResponseFrame(f xbee.APIFrame) error {
	var frame xbee.RemoteATCommandResponseAPIFrameData

	//TODO move that in xbee/api_frame_data_remote_at_command_request.go in a FromBytes functions?
	offset := 0
	frame.Type = f.Data[offset]
	offset += 1
	frame.ID = f.Data[offset]
	offset += 1
	copy(frame.SourceAddress64[:], f.Data[offset:offset+8])
	offset += 8
	copy(frame.SourceAddress16[:], f.Data[offset:offset+2])
	offset += 2
	copy(frame.ATCommand[:], f.Data[offset:offset+2])
	offset += 2
	frame.CommandStatus = f.Data[offset]
	offset += 1

	frame.ParameterValue = make([]byte, len(f.Data[offset:]))
	copy(frame.ParameterValue[:], f.Data[offset:])

	if frame.CommandStatus != 0x00 {
		var msg string
		switch frame.CommandStatus {
		case 0x01:
			msg = "ERROR"
		case 0x02:
			msg = "Invalid command"
		case 0x03:
			msg = "Invalid parameter"
		case 0x04:
			msg = "Transmission failure"
		default:
			msg = "Unknwon status code"
		}

		return fmt.Errorf("remote_at_command_response status is not OK (%X=%s)", frame.CommandStatus, msg)
	}

	switch strings.ToUpper(string(frame.ATCommand[:])) {
	case "SH":
		xbee.RecordSH(hex.EncodeToString(frame.SourceAddress16[:]), hex.EncodeToString(frame.ParameterValue[:]))
		if err := mqtt.Publish("xbee/remote_at_commmand_response", frame); err != nil {
			log.Errorln(err)
		}
	case "SL":
		xbee.RecordSL(hex.EncodeToString(frame.SourceAddress16[:]), hex.EncodeToString(frame.ParameterValue[:]))
	}

	return nil
}
