package xbee

import (
	"bytes"
	"encoding/binary"
)

func WriteAPIFameDataToSerial(fd APIFrameData) error {
	ba, _ := fd.Bytes()

	var buf bytes.Buffer
	buf.WriteByte(0x7e)                                   // Start Delimiter
	binary.Write(&buf, binary.BigEndian, uint16(len(ba))) // Length
	buf.Write(ba)                                         // Frame Data
	buf.WriteByte(computeAPIFrameChecksum(ba))            // Checksum

	_, err := serialPort.Write(buf.Bytes())

	return err
}

func WriteSerial(f APIFrame) error {
	var buf bytes.Buffer
	buf.WriteByte(f.StartDelimiter)
	binary.Write(&buf, binary.BigEndian, f.Length) // buf.Write(f.Length[:])
	buf.Write(f.Data)
	buf.WriteByte(f.Checksum)

	_, err := serialPort.Write(buf.Bytes())

	return err
}
