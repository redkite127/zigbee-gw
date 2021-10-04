// Transmit Request - 0x10
// https://www.digi.com/resources/documentation/DigiDocs/90002002/Default.htm#Reference/r_frame_0x10.htm
package xbee

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const transmitRequestFrameDataMinimumSize = 14

type TransmitRequestFrameData struct {
	Type                 byte
	ID                   uint8 // byte
	DestinationAddress64 [8]byte
	DestinationAddress16 [2]byte
	BroadcastRadius      byte
	TransmitOptions      byte
	PayloadData          []byte
}

func NewTransmitRequestFrameData() TransmitRequestFrameData {
	return TransmitRequestFrameData{
		Type:                 0x10,
		ID:                   0x01,
		DestinationAddress64: [8]byte{},
		DestinationAddress16: [2]byte{},
		BroadcastRadius:      0x00,
		TransmitOptions:      0x00,
		PayloadData:          []byte{},
	}
}

func (f *TransmitRequestFrameData) SetDestinationAddress64(address64 string) error {
	dest64ba, err := hex.DecodeString(address64)
	if err != nil {
		return err
	}

	copy(f.DestinationAddress64[:], dest64ba)
	copy(f.DestinationAddress16[:], []byte{0xff, 0xfe})

	return nil
}

func (f *TransmitRequestFrameData) SetDestinationAddress16(address16 string) error {
	dest16ba, err := hex.DecodeString(address16)
	if err != nil {
		return err
	}

	copy(f.DestinationAddress16[:], dest16ba)
	copy(f.DestinationAddress64[:], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	return nil
}

// func (f *TransmitRequestFrameData) SetBroadcastRadius(radius uint8) error

// func (f *TransmitRequestFrameData) SetTransmitOptions(options byte) error

// func (f *TransmitRequestFrameData) SetPayloadData(payload []byte) error

// MarshalBinary implements encoding.BinaryMarshaler.
func (f TransmitRequestFrameData) MarshalBinary() ([]byte, error) {
	//TODO how to reduce memory consumption and return immediatelly a byte array or a reader?
	var buf bytes.Buffer
	buf.Grow(transmitRequestFrameDataMinimumSize)

	buf.WriteByte(f.Type)
	binary.Write(&buf, binary.BigEndian, f.ID)
	buf.Write(f.DestinationAddress64[:])
	buf.Write(f.DestinationAddress16[:])
	buf.WriteByte(f.BroadcastRadius)
	buf.WriteByte(f.TransmitOptions)
	buf.Write(f.PayloadData)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (f *TransmitRequestFrameData) UnmarshalBinary(data []byte) error {
	if len(data) < transmitRequestFrameDataMinimumSize {
		return fmt.Errorf("invalid 'Transmit Request' frame data (got only %d bytes)", len(data))
	}
	if TransmitAPIFrameType(data[0]) != TypeTransmitRequest {
		return fmt.Errorf("invalid 'Transmit Request' frame data type (got %X instead of %X)", data[0], TypeTransmitRequest)
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
	f.BroadcastRadius = data[offset]
	offset += 1
	f.TransmitOptions = data[offset]
	offset += 1
	f.PayloadData = make([]byte, len(data[offset:]))
	copy(f.PayloadData[:], data[offset:])

	return nil
}
