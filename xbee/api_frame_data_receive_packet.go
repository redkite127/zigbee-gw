// Receive Packet - 0x90
// https://www.digi.com/resources/documentation/Digidocs/90002002/#Reference/r_frame_0x90.htm
package xbee

import (
	"bytes"
)

type ReceivePacketAPIFrameData struct {
	Type            byte
	SourceAddress64 [8]byte
	SourceAddress16 [2]byte
	ReceiveOptions  byte
	ReceivedData    []byte
}

//TODO really useful for frame that we receive? Don't we only need this for transmit frame types?
func NewReceivePacketAPIFrameData() ReceivePacketAPIFrameData {
	return ReceivePacketAPIFrameData{
		Type:            0x90,
		SourceAddress64: [8]byte{},
		SourceAddress16: [2]byte{},
		ReceiveOptions:  0x01,
		ReceivedData:    []byte{},
	}
}

//TODO really useful for frame that we receive? Don't we only need this for transmit frame types?
func (f ReceivePacketAPIFrameData) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteByte(f.Type)
	buf.Write(f.SourceAddress64[:])
	buf.Write(f.SourceAddress16[:])
	buf.WriteByte(f.ReceiveOptions)
	buf.Write(f.ReceivedData)

	return buf.Bytes()
}
