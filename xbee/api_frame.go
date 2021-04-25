// API frame structure
// https://www.digi.com/resources/documentation/Digidocs/90001942-13/#concepts/c_api_frame_structure.htm
package xbee

import "encoding/binary"

type apiFrameSegment uint8

const (
	segmentStartDelimiter = apiFrameSegment(iota)
	segmentLength         = apiFrameSegment(iota + 1)
	segmentData           = apiFrameSegment(iota + 2)
	segmentChecksum       = apiFrameSegment(iota + 3)
)

type APIFrame struct {
	StartDelimiter byte
	Length         uint16 // [2]byte
	Data           []byte
	Checksum       byte
}

func computeAPIFrameChecksum(ba []byte) byte {
	var sum uint64
	for _, b := range ba {
		sum += uint64(b)
	}

	var sumb [8]byte
	binary.BigEndian.PutUint64(sumb[:], sum)
	lastByte := sumb[len(sumb)-1]
	return 0xff - lastByte
}
