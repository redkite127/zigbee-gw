// Remote AT Command Response - 0x97
// https://www.digi.com/resources/documentation/Digidocs/90002002/#Reference/r_frame_0x97.htm
package xbee

import (
	"bytes"
	"encoding/binary"
)

type RemoteATCommandResponseAPIFrameData struct {
	Type            byte
	ID              uint8 // byte
	SourceAddress64 [8]byte
	SourceAddress16 [2]byte
	ATCommand       [2]byte
	CommandStatus   byte
	ParameterValue  []byte
}

func (f RemoteATCommandResponseAPIFrameData) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteByte(f.Type)
	binary.Write(&buf, binary.BigEndian, f.ID)
	buf.Write(f.SourceAddress64[:])
	buf.Write(f.SourceAddress16[:])
	buf.Write(f.ATCommand[:])
	buf.WriteByte(f.CommandStatus)
	buf.Write(f.ParameterValue)

	return buf.Bytes()
}
