package xbee

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var serialConfig *serial.Config
var serialPort *serial.Port

func InitSerial(name string, speed int) {
	serialConfig = &serial.Config{Name: name, Baud: speed /*, ReadTimeout: 5 * time.Second*/}

	var err error
	if serialPort, err = serial.OpenPort(serialConfig); err != nil {
		log.Fatal(err)
	}
}

func SendRemoteATcommand(dest64, dest16, atCommand string) error {
	// handle destination address
	if dest64 == "" && dest16 != "" {
		dest64 = "FFFFFFFFFFFFFFFF"
	} else if dest16 == "" && dest64 != "" {
		dest16 = "FFFE"
	} else if dest16 != "" && dest64 != "" {
		dest16 = "FFFE"
	} else {
		return fmt.Errorf("no valid destination address given")
	}

	data := make([]byte, 15)

	data[0] = 0x17                         // frame type
	data[1] = 0x01                         // frame id
	dest64b, _ := hex.DecodeString(dest64) // prepare 64-bits
	copy(data[2:10], dest64b)              // 64-bits destination address
	dest16b, _ := hex.DecodeString(dest16) // prepare 16-bits
	copy(data[10:12], dest16b)             // 16-bits destination address
	data[12] = 0x00                        // command options
	copy(data[13:15], []byte(atCommand))   // AT command

	//log.Debugf("% X", data)

	return SendData(data)
}

func SendData(data []byte) error {
	var f Frame
	f.StartDelimiter = 0x7E
	binary.BigEndian.PutUint16(f.Length[:], uint16(len(data)))
	f.Data = data

	// compute checksum
	var sum uint64
	for _, b := range data {
		sum += uint64(b)
	}
	var sumb [8]byte
	binary.BigEndian.PutUint64(sumb[:], sum)
	lastByte := sumb[len(sumb)-1]
	f.Checksum = 0xff - lastByte

	return sendFrame(f)
}

func sendFrame(f Frame) error {
	var buf bytes.Buffer
	buf.Write([]byte{f.StartDelimiter})
	buf.Write(f.Length[:])
	buf.Write(f.Data)
	buf.Write([]byte{f.Checksum})

	return sendSerial(buf.Bytes())
}

func sendSerial(b []byte) error {
	_, err := serialPort.Write(b)

	return err
}

func ReadSerial(fc chan<- ReceivePacketFrame) {
	var frame ReceivePacketFrame
	var escaping bool
	var buffer bytes.Buffer
	buf := make([]byte, 1)
	for {
		n, err := serialPort.Read(buf)
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
