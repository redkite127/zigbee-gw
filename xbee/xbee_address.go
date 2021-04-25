package xbee

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

//shortToLongAddress will be filled and updated with packet received
var shortToLongAddress = map[string]string{}
var shortToSH = map[string]string{}
var shortToSL = map[string]string{}

func Record16bitsAddress(address64, address16 string) {
	shortToLongAddress[address16] = address64
	log.Debugln(shortToLongAddress)
}

func RecordSH(address16, sh string) {
	shortToSH[address16] = sh
	log.Debugln("SH:", shortToSH)
}

func RecordSL(address16, sl string) {
	shortToSL[address16] = sl
	log.Debugln("SL:", shortToSL)
}

func Fix64address(address16 string) (string, error) {
	fixed64, found := shortToLongAddress[address16]
	if !found {
		log.Debugln("no match found for the 16-bits address")
		//ask SH & SL to the device using the Remote AT commands
		return Get64addressFrom16address(address16)
	}

	return fixed64, nil
}

func Get64addressFrom16address(address16 string) (string, error) {
	r1 := NewRemoteATCommandRequestFrameData()
	r1.SetDestinationAddress16(address16)
	r1.SetATCommand("SH")
	WriteAPIFameDataToSerial(&r1)

	r2 := NewRemoteATCommandRequestFrameData()
	r2.SetDestinationAddress16(address16)
	r2.SetATCommand("SL")
	WriteAPIFameDataToSerial(&r2)

	time.Sleep(5 * time.Second)

	sh, found1 := shortToSH[address16]
	sl, found2 := shortToSL[address16]
	if !found1 || !found2 {
		return "", fmt.Errorf("failed to retrieve SH & SL (%t;%t)", found1, found2)
	}

	Record16bitsAddress(sh+sl, address16)

	delete(shortToSH, address16)
	delete(shortToSL, address16)

	return sh + sl, nil
}
