package xbee

type APIFrameData interface {
	Bytes() ([]byte, error)
	FromBytes(b []byte) error
}
