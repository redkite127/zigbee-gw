package xbee

import (
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

func ReadSerial(fc chan<- APIFrame) {
	log := log.WithField("backend", "xbee")

	if serialPort == nil {
		log.Fatal("serial hasn't been initialized (xbee.InitSerial)")
	}

	var frame APIFrame
	var frameSegment apiFrameSegment = segmentStartDelimiter
	var escaping bool
	var buffer bytes.Buffer
	buf := make([]byte, 1)
	for {
		// read one byte
		n, err := serialPort.Read(buf)
		if err != nil {
			log.Fatalf("failed to read from serial port: %w", err)
		} else if n <= 0 {
			//log.Debug("haven't read any byte")
			continue
		}
		b := buf[0]

		// Escape Byte
		// Read the next character and don't count this one.
		if b == 0x7D {
			escaping = true
			continue
		} else if escaping {
			b ^= 0x20
			escaping = false
		}
		//log.Debugf("%X", b)

		// Start Delimiter
		// No matter the current state, we start a new frame.
		if b == 0x7E {
			frame = APIFrame{StartDelimiter: b}
			frameSegment = segmentLength
			escaping = false
			buffer.Reset()
			continue
		}

		switch frameSegment {
		case segmentLength:
			buffer.WriteByte(b)
			if buffer.Len() == 2 {
				frame.Length = binary.BigEndian.Uint16(buffer.Bytes())
				frameSegment = segmentData
				buffer.Reset()
			}
		case segmentData:
			buffer.WriteByte(b)
			if buffer.Len() == int(frame.Length) {
				frame.Data = buffer.Bytes()
				frameSegment = segmentChecksum
				buffer.Reset()
			}
		case segmentChecksum:
			frame.Checksum = b
			frameSegment = segmentStartDelimiter

			if cc := computeAPIFrameChecksum(frame.Data); cc != frame.Checksum {
				log.Errorf("discarding the frame because computed checksum is not matching the expected checksum: %X != %X", cc, frame.Checksum)
				continue
			}

			fc <- frame
		}
	}
}
