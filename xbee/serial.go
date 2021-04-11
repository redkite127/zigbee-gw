package xbee

import (
	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

var serialConfig *serial.Config
var serialPort *serial.Port

func InitSerial(name string, speed int) {
	serialConfig = &serial.Config{Name: name, Baud: speed /*, ReadTimeout: 5 * time.Second*/}

	var err error
	if serialPort, err = serial.OpenPort(serialConfig); err != nil {
		log.Fatal(err)
	}
}
