package main

import "github.com/redkite1/zigbee-gw/xbee"

func TCPredirect(destination, data string) error {
	req := xbee.NewTransmitRequestFrameData()
	if len(destination) > 4 {
		req.SetDestinationAddress64(destination)
	} else {
		req.SetDestinationAddress16(destination)
	}
	req.SetPayloadData([]byte(data))

	return xbee.WriteAPIFameDataToSerial(&req)
}
