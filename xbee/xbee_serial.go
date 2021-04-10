package xbee

import (
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

func ReadSerial(fc chan<- ReceivePacketFrame, name string, speed int) {
	c := &serial.Config{Name: name, Baud: speed /*, ReadTimeout: 5 * time.Second*/}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	var frame ReceivePacketFrame
	var escaping bool
	var buffer bytes.Buffer
	buf := make([]byte, 1)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		} else if n <= 0 {
			log.Debugf("Haven't read any byte")
			continue
		}
		//log.Printf("Received %d bytes: %x", n, buf)

		// Handle start byte
		b := buf[0]
		if b == 0x7E {
			//log.Debugf("FrameStart received!")
			frame.State = StateLength
			escaping = false
			buffer.Reset()
			continue
		}

		// Handle escape byte
		if b == 0x7D {
			//log.Debugf("Escape character received!")
			escaping = true
			continue
		} else if escaping {
			escaping = false
			b ^= 0x20
		}

		//log.Printf("%X", b)
		switch frame.State {
		case StateLength:
			buffer.WriteByte(b)
			if buffer.Len() == 2 {
				frame.Length = binary.BigEndian.Uint16(buffer.Next(2))
				//log.Debugf("Frame length: ", frame.Length)
				frame.State = StateType
			}
		case StateType:
			frame.Type = FrameType(b)
			//log.Debugf("Frame type: %X\n", frame.Type)
			frame.State = StateData
		case StateData:
			buffer.WriteByte(b)
			if buffer.Len() == int(frame.Length)-1 {
				frame.Data = buffer.Next(int(frame.Length) - 1)
				//log.Debugf("Frame data: % X\n", frame.Data)
				frame.State = StateChecksum
			}
		case StateChecksum:
			frame.Checksum = b
			//log.Debugf("Frame checksum: %X\n", frame.Checksum)
			fc <- frame
			frame.State = StateStart
		}
	}
}
