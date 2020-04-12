package xbee

type Frame struct {
	State    FrameState
	Length   uint16
	Type     FrameType
	Data     []byte
	Checksum byte
}

type FrameState uint8

const (
	StateStart    = FrameState(iota)
	StateLength   = FrameState(iota + 1)
	StateType     = FrameState(iota + 2)
	StateData     = FrameState(iota + 3)
	StateChecksum = FrameState(iota + 4)
)

type FrameType uint8

const (
	// https://www.digi.com/resources/documentation/Digidocs/90001942-13/#reference/r_supported_frames_zigbee.htm?Highlight=receive packet
	TypeATCommandResponse              = FrameType(0x88)
	TypeModemStatus                    = FrameType(0x8A)
	TypeTransmitStatus                 = FrameType(0x88)
	TypeReceivePacket                  = FrameType(0x90)
	TypeExplicitRxIndicator            = FrameType(0x91)
	TypeIODataSampleRxIndicator        = FrameType(0x92)
	TypeXBeeSensorReadIndicator        = FrameType(0x94)
	TypeNodeIdentificationIndicator    = FrameType(0x95)
	TypeRemoteATCommandResponse        = FrameType(0x97)
	TypeExtendedModemStatus            = FrameType(0x98)
	TypeOverTheAirFirmwareUpdateStatus = FrameType(0xA0)
	TypeRouterRecordIndicator          = FrameType(0xA1)
	TypeManyToOneRouteRequestIndicator = FrameType(0xA3)
	TypeJoinNotificationStatus         = FrameType(0xA5)
)
