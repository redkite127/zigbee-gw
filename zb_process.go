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
			//err = processReceivePacketFrame(f)
			go processReceivePacketFrame(f)
		case xbee.TypeRemoteATCommandResponse:
			go processRemoteATCommandResponseFrame(f)
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
	// frame := xbee.ReceivePacketFrameData{}
	// buf := bytes.NewReader(f.Data)
	// err := binary.Read(buf, binary.BigEndian, &frame)
	// if err != nil {
	// 	log.Errorf("failed to decode 'Receive Packet' frame data: %w", err)
	// 	return fmt.Errorf("failed to decode 'Receive Packet' frame data: %w", err)
	// }
	//WHY UnmarshalBinary is not called?

	var frame xbee.ReceivePacketFrameData
	if err := frame.FromBytes(f.Data); err != nil {
		log.Errorf("failed to decode 'Receive Packet' frame data: %s", err)
		return err
	}

	log.Debugf("64-bit source address: % X", frame.SourceAddress64)
	log.Debugf("16-bit source address: % X", frame.SourceAddress16)
	log.Debugf("Receive options: % X", frame.ReceiveOptions)
	log.Debugf("RF data: % X", frame.ReceivedData)
	log.Debugf("RF data (string): %s", frame.ReceivedData)

	if err := ZBredirect(hex.EncodeToString(frame.SourceAddress64[:]), hex.EncodeToString(frame.SourceAddress16[:]), frame.ReceivedData); err != nil {
		return err
	}

	return nil
}

func processRemoteATCommandResponseFrame(f xbee.APIFrame) error {
	var frame xbee.RemoteATCommandResponseFrameData
	if err := frame.FromBytes(f.Data); err != nil {
		log.Errorf("failed to decode 'Remote AT Command Response' frame data: %s", err)
		return err
	}

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

		return fmt.Errorf("'Remote AT Command Response' status is not OK (%X = %s)", frame.CommandStatus, msg)
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
