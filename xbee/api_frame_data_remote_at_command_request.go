// Remote AT Command Request - 0x17
// https://www.digi.com/resources/documentation/Digidocs/90002002/#Reference/r_frame_0x17.htm
package xbee

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const remoteATCommandRequestFrameDataMinimumSize = 15

type RemoteATCommandRequestFrameData struct {
	Type                 byte
	ID                   uint8 // byte
	DestinationAddress64 [8]byte
	DestinationAddress16 [2]byte
	RemoteCommandOptions byte
	ATCommand            [2]byte
	ParameterValue       []byte
}

func NewRemoteATCommandRequestFrameData() RemoteATCommandRequestFrameData {
	return RemoteATCommandRequestFrameData{
		Type:                 0x17,
		ID:                   0x01,
		DestinationAddress64: [8]byte{},
		DestinationAddress16: [2]byte{},
		RemoteCommandOptions: 0x02,
		ATCommand:            [2]byte{},
		ParameterValue:       []byte{},
	}
}

func (f *RemoteATCommandRequestFrameData) SetDestinationAddress64(address64 string) error {
	dest64ba, err := hex.DecodeString(address64)
	if err != nil {
		return err
	}

	copy(f.DestinationAddress64[:], dest64ba)
	copy(f.DestinationAddress16[:], []byte{0xff, 0xfe})

	return nil
}

func (f *RemoteATCommandRequestFrameData) SetDestinationAddress16(address16 string) error {
	dest16ba, err := hex.DecodeString(address16)
	if err != nil {
		return err
	}

	copy(f.DestinationAddress16[:], dest16ba)
	copy(f.DestinationAddress64[:], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	return nil
}

func (f *RemoteATCommandRequestFrameData) SetATCommand(ATCommand string) error {
	if len(ATCommand) != 2 {
		return fmt.Errorf("invalid AT command")
	}

	copy(f.ATCommand[:], []byte(ATCommand))

	return nil
}

func (f *RemoteATCommandRequestFrameData) FromBytes(b []byte) error {
	return f.UnmarshalBinary(b)
}

//TODO really useful for frame that we receive? Don't we only need this for transmit frame types?
func (f RemoteATCommandRequestFrameData) Bytes() ([]byte, error) {
	return f.MarshalBinary()
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (f RemoteATCommandRequestFrameData) MarshalBinary() ([]byte, error) {
	//TODO how to reduce memory consumption and return immediatelly a byte array or a reader?
	var buf bytes.Buffer
	buf.Grow(remoteATCommandRequestFrameDataMinimumSize)

	buf.WriteByte(f.Type)
	binary.Write(&buf, binary.BigEndian, f.ID)
	buf.Write(f.DestinationAddress64[:])
	buf.Write(f.DestinationAddress16[:])
	buf.WriteByte(f.RemoteCommandOptions)
	buf.Write(f.ATCommand[:])
	buf.Write(f.ParameterValue)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (f *RemoteATCommandRequestFrameData) UnmarshalBinary(data []byte) error {
	if len(data) < remoteATCommandRequestFrameDataMinimumSize {
		return fmt.Errorf("invalid 'Remote AT Command Request' frame data (got only %d bytes)", len(data))
	}
	if TransmitAPIFrameType(data[0]) != TypeRemoteATCommandRequest {
		return fmt.Errorf("invalid 'Remote AT Command Request' frame data type (got %X instead of %X)", data[0], TypeRemoteATCommandRequest)
	}

	offset := 0
	f.Type = data[offset]
	offset += 1
	f.ID = data[offset]
	offset += 1
	copy(f.DestinationAddress64[:], data[offset:offset+8])
	offset += 8
	copy(f.DestinationAddress16[:], data[offset:offset+2])
	offset += 2
	f.RemoteCommandOptions = data[offset]
	offset += 1
	copy(f.ATCommand[:], data[offset:offset+2])
	offset += 2
	f.ParameterValue = make([]byte, len(data[offset:]))
	copy(f.ParameterValue[:], data[offset:])

	return nil
}
