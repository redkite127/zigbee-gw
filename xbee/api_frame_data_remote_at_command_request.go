// Remote AT Command Request - 0x17
// https://www.digi.com/resources/documentation/Digidocs/90002002/#Reference/r_frame_0x17.htm
package xbee

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type RemoteATCommandRequestAPIFrameData struct {
	Type                 byte
	ID                   uint8 // byte
	DestinationAddress64 [8]byte
	DestinationAddress16 [2]byte
	RemoteCommandOptions byte
	ATCommand            [2]byte
	ParameterValue       []byte
}

func NewRemoteATCommandRequestAPIFrameData() RemoteATCommandRequestAPIFrameData {
	return RemoteATCommandRequestAPIFrameData{
		Type:                 0x17,
		ID:                   0x01,
		DestinationAddress64: [8]byte{},
		DestinationAddress16: [2]byte{},
		RemoteCommandOptions: 0x02,
		ATCommand:            [2]byte{},
		ParameterValue:       []byte{},
	}
}

func (f RemoteATCommandRequestAPIFrameData) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteByte(f.Type)
	binary.Write(&buf, binary.BigEndian, f.ID)
	buf.Write(f.DestinationAddress64[:])
	buf.Write(f.DestinationAddress16[:])
	buf.WriteByte(f.RemoteCommandOptions)
	buf.Write(f.ATCommand[:])
	buf.Write(f.ParameterValue)

	return buf.Bytes()
}

func (f *RemoteATCommandRequestAPIFrameData) SetDestinationAddress64(address64 string) error {
	dest64ba, err := hex.DecodeString(address64)
	if err != nil {
		return err
	}

	copy(f.DestinationAddress64[:], dest64ba)
	copy(f.DestinationAddress16[:], []byte{0xff, 0xfe})

	return nil
}

func (f *RemoteATCommandRequestAPIFrameData) SetDestinationAddress16(address16 string) error {
	dest16ba, err := hex.DecodeString(address16)
	if err != nil {
		return err
	}

	copy(f.DestinationAddress16[:], dest16ba)
	copy(f.DestinationAddress64[:], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	return nil
}

func (f *RemoteATCommandRequestAPIFrameData) SetATCommand(ATCommand string) error {
	if len(ATCommand) != 2 {
		return fmt.Errorf("invalid AT command")
	}

	copy(f.ATCommand[:], []byte(ATCommand))

	return nil
}
