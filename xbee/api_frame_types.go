// Supported frames
// https://www.digi.com/resources/documentation/Digidocs/90001942-13/#reference/r_supported_frames_zigbee.htm
package xbee

type ReceiveAPIFrameType uint8

const (
	TypeATCommandResponse              = ReceiveAPIFrameType(0x88)
	TypeModemStatus                    = ReceiveAPIFrameType(0x8A)
	TypeTransmitStatus                 = ReceiveAPIFrameType(0x8B)
	TypeReceivePacket                  = ReceiveAPIFrameType(0x90)
	TypeExplicitRxIndicator            = ReceiveAPIFrameType(0x91)
	TypeIODataSampleRxIndicator        = ReceiveAPIFrameType(0x92)
	TypeXBeeSensorReadIndicator        = ReceiveAPIFrameType(0x94)
	TypeNodeIdentificationIndicator    = ReceiveAPIFrameType(0x95)
	TypeRemoteATCommandResponse        = ReceiveAPIFrameType(0x97)
	TypeExtendedModemStatus            = ReceiveAPIFrameType(0x98)
	TypeOverTheAirFirmwareUpdateStatus = ReceiveAPIFrameType(0xA0)
	TypeRouterRecordIndicator          = ReceiveAPIFrameType(0xA1)
	TypeManyToOneRouteRequestIndicator = ReceiveAPIFrameType(0xA3)
	TypeJoinNotificationStatus         = ReceiveAPIFrameType(0xA5)
)

type TransmitAPIFrameType uint8

const (
	TypeATCommand                      = TransmitAPIFrameType(0x08)
	TypeATCommandQueueParameterValue   = TransmitAPIFrameType(0x09)
	TypeTransmitRequest                = TransmitAPIFrameType(0x10)
	TypeExplicitAddressingCommandFrame = TransmitAPIFrameType(0x11)
	TypeRemoteATCommandRequest         = TransmitAPIFrameType(0x17)
	TypeCreateSourceRoute              = TransmitAPIFrameType(0x21)
	TypeRegisterJoiningDevice          = TransmitAPIFrameType(0x24)
)
