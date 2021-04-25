// Remote AT Command Response - 0x97
// https://www.digi.com/resources/documentation/Digidocs/90002002/#Reference/r_frame_0x97.htm
package xbee

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const remoteATCommandResponseFrameDataMinimumSize = 15

type RemoteATCommandResponseFrameData struct {
	Type            byte
	ID              uint8 // byte
	SourceAddress64 [8]byte
	SourceAddress16 [2]byte
	ATCommand       [2]byte
	CommandStatus   byte
	ParameterValue  []byte
}

func (f *RemoteATCommandResponseFrameData) FromBytes(b []byte) error {
	return f.UnmarshalBinary(b)
}

//TODO really useful for frame that we receive? Don't we only need this for transmit frame types?
func (f RemoteATCommandResponseFrameData) Bytes() ([]byte, error) {
	return f.MarshalBinary()
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (f RemoteATCommandResponseFrameData) MarshalBinary() ([]byte, error) {
	//TODO how to reduce memory consumption and return immediatelly a byte array or a reader?
	var buf bytes.Buffer
	buf.Grow(remoteATCommandResponseFrameDataMinimumSize)

	buf.WriteByte(f.Type)
	binary.Write(&buf, binary.BigEndian, f.ID)
	buf.Write(f.SourceAddress64[:])
	buf.Write(f.SourceAddress16[:])
	buf.Write(f.ATCommand[:])
	buf.WriteByte(f.CommandStatus)
	buf.Write(f.ParameterValue)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (f *RemoteATCommandResponseFrameData) UnmarshalBinary(data []byte) error {
	if len(data) < remoteATCommandResponseFrameDataMinimumSize {
		return fmt.Errorf("invalid 'Remote AT Command Response' frame data (got only %d bytes)", len(data))
	}
	if ReceiveAPIFrameType(data[0]) != TypeRemoteATCommandResponse {
		return fmt.Errorf("invalid 'Remote AT Command Response' frame data type (got %X instead of %X)", data[0], TypeRemoteATCommandResponse)
	}

	offset := 0
	f.Type = data[offset]
	offset += 1
	f.ID = data[offset]
	offset += 1
	copy(f.SourceAddress64[:], data[offset:offset+8])
	offset += 8
	copy(f.SourceAddress16[:], data[offset:offset+2])
	offset += 2
	copy(f.ATCommand[:], data[offset:offset+2])
	offset += 2
	f.CommandStatus = data[offset]
	offset += 1
	f.ParameterValue = make([]byte, len(data[offset:]))
	copy(f.ParameterValue[:], data[offset:])

	return nil
}
