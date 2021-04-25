// Receive Packet - 0x90
// https://www.digi.com/resources/documentation/Digidocs/90002002/#Reference/r_frame_0x90.htm
package xbee

import (
	"bytes"
	"fmt"
)

const receivePacketFrameDataMinimumSize = 12

type ReceivePacketFrameData struct {
	Type            byte
	SourceAddress64 [8]byte
	SourceAddress16 [2]byte
	ReceiveOptions  byte
	ReceivedData    []byte
}

//TODO really useful for frame that we receive? Don't we only need this for transmit frame types?
func NewReceivePacketFrameData() ReceivePacketFrameData {
	return ReceivePacketFrameData{
		Type:            0x90,
		SourceAddress64: [8]byte{},
		SourceAddress16: [2]byte{},
		ReceiveOptions:  0x01,
		ReceivedData:    []byte{},
	}
}

func (f *ReceivePacketFrameData) FromBytes(b []byte) error {
	return f.UnmarshalBinary(b)
}

//TODO really useful for frame that we receive? Don't we only need this for transmit frame types?
func (f ReceivePacketFrameData) Bytes() ([]byte, error) {
	return f.MarshalBinary()
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (f ReceivePacketFrameData) MarshalBinary() ([]byte, error) {
	//TODO how to reduce memory consumption and return immediatelly a byte array or a reader?
	var buf bytes.Buffer
	buf.Grow(receivePacketFrameDataMinimumSize)

	buf.WriteByte(f.Type)
	buf.Write(f.SourceAddress64[:])
	buf.Write(f.SourceAddress16[:])
	buf.WriteByte(f.ReceiveOptions)
	buf.Write(f.ReceivedData)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (f *ReceivePacketFrameData) UnmarshalBinary(data []byte) error {
	if len(data) < receivePacketFrameDataMinimumSize {
		return fmt.Errorf("invalid 'Receive Packet' frame data (got only %d bytes)", len(data))
	}
	if ReceiveAPIFrameType(data[0]) != TypeReceivePacket {
		return fmt.Errorf("invalid 'Receive Packet' frame data type (got %X instead of %X)", data[0], TypeReceivePacket)
	}

	offset := 0
	f.Type = data[offset]
	offset += 1
	copy(f.SourceAddress64[:], data[offset:offset+8])
	offset += 8
	copy(f.SourceAddress16[:], data[offset:offset+2])
	offset += 2
	f.ReceiveOptions = data[offset]
	offset += 1
	f.ReceivedData = make([]byte, len(data[offset:]))
	copy(f.ReceivedData[:], data[offset:])

	return nil
}
