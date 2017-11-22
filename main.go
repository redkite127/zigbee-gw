package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/tarm/serial"
)

type FrameStates int

const (
	FrameStart    = FrameStates(0)
	FrameLength   = FrameStart + 1
	FrameType     = FrameLength + 1
	FrameData     = FrameType + 1
	FrameChecksum = FrameData + 1
)

type FrameTypes uint8

const (
	// https://www.digi.com/resources/documentation/Digidocs/90001942-13/#reference/r_supported_frames_zigbee.htm?Highlight=receive packet
	FrameTypeATCommandResponse              = FrameTypes(0x88)
	FrameTypeModemStatus                    = FrameTypes(0x8A)
	FrameTypeTransmitStatus                 = FrameTypes(0x88)
	FrameTypeReceivePacket                  = FrameTypes(0x90)
	FrameTypeExplicitRxIndicator            = FrameTypes(0x91)
	FrameTypeIODataSampleRxIndicator        = FrameTypes(0x92)
	FrameTypeXBeeSensorReadIndicator        = FrameTypes(0x94)
	FrameTypeNodeIdentificationIndicator    = FrameTypes(0x95)
	FrameTypeRemoteATCommandResponse        = FrameTypes(0x97)
	FrameTypeExtendedModemStatus            = FrameTypes(0x98)
	FrameTypeOverTheAirFirmwareUpdateStatus = FrameTypes(0xA0)
	FrameTypeRouterRecordIndicator          = FrameTypes(0xA1)
	FrameTypeManyToOneRouteRequestIndicator = FrameTypes(0xA3)
	FrameTypeJoinNotificationStatus         = FrameTypes(0xA5)
)

type XBeeFrame struct {
	State    FrameStates
	Length   uint16
	Type     FrameTypes
	Data     []byte
	Checksum byte
}

type Redirection struct {
	Scheme     string
	Host       string
	Port       int
	Path       string
	Parameters url.Values
}

var registeredRedirection map[string]Redirection = map[string]Redirection{
	"0013a20041531c31": {
		Scheme: "http",
		Host:   "10.161.0.130",
		Port:   2001,
		Path:   "sensors",
		Parameters: url.Values{
			//"name": {"cave_temperature_humidity_1"},
			"room": {"cave"},
			"type": {"temperature;humidity"},
		},
	},
}

func readSerial(fc chan<- XBeeFrame) {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600 /*, ReadTimeout: 5 * time.Second*/}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	var frame XBeeFrame
	var escaping bool
	var buffer bytes.Buffer
	buf := make([]byte, 1)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		} else if n <= 0 {
			log.Println("Haven't read any byte")
			continue
		}
		//log.Printf("Received %d bytes: %x", n, buf)

		// Handle start byte
		b := buf[0]
		if b == 0x7E {
			//log.Println("FrameStart received!")
			frame.State = FrameLength
			escaping = false
			buffer.Reset()
			continue
		}

		// Handle escape byte
		if b == 0x7D {
			//log.Println("Escape character received!")
			escaping = true
			continue
		} else if escaping {
			escaping = false
			b ^= 0x20
		}

		//log.Printf("%X", b)
		switch frame.State {
		case FrameLength:
			buffer.WriteByte(b)
			if buffer.Len() == 2 {
				frame.Length = binary.BigEndian.Uint16(buffer.Next(2))
				//log.Println("Frame length: ", frame.Length)
				frame.State = FrameType
			}
		case FrameType:
			frame.Type = FrameTypes(b)
			//log.Printf("Frame type: %X\n", frame.Type)
			frame.State = FrameData
		case FrameData:
			buffer.WriteByte(b)
			if buffer.Len() == int(frame.Length)-1 {
				frame.Data = buffer.Next(int(frame.Length) - 1)
				//log.Printf("Frame data: % X\n", frame.Data)
				frame.State = FrameChecksum
			}
		case FrameChecksum:
			frame.Checksum = b
			//log.Printf("Frame checksum: %X\n", frame.Checksum)

			fc <- frame
		}
	}
}

func processFrames(fc <-chan XBeeFrame) {
	var err error
	for f := range fc {
		switch f.Type {
		case FrameTypeReceivePacket:
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

func processReceivePacketFrame(f XBeeFrame) error {
	sa64 := f.Data[:8]
	log.Printf("64-bit source address: % X", sa64)

	sa16 := f.Data[8:10]
	log.Printf("16-bit source address: % X", sa16)

	ro := f.Data[10:11]
	log.Printf("Receive options: % X", ro)

	rfd := f.Data[11:]
	log.Printf("RF data: % X", rfd)
	log.Printf("RF data (string): %s", rfd)

	r, found := registeredRedirection[hex.EncodeToString(sa64)]
	if !found {
		return errors.New("Received packet from unregistered device!")
	}

	u := url.URL{
		Scheme:   r.Scheme,
		Host:     r.Host + ":" + strconv.Itoa(r.Port),
		Path:     r.Path,
		RawQuery: r.Parameters.Encode(),
	}

	if resp, err := http.Post(u.String(), "text/plain", bytes.NewReader(rfd)); err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg string
		if err != nil {
			msg = err.Error()
		} else {
			msg = resp.Status
		}

		return errors.New("Failed to redirect data to destination: " + msg)
	}

	return nil
}

func main() {
	//stopChan := make(chan interface{})
	frameChan := make(chan XBeeFrame)
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		close(frameChan)
	}()

	go readSerial(frameChan /*, stopChan*/)
	processFrames(frameChan)
}
