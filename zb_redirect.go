package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/redkite1/zigbee-gw/xbee"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ZBSource struct {
	Room   string
	Format string
}

var registeredZBSources map[string]ZBSource

func ZBredirect(source64, source16 string, data []byte) (err error) {
	// looking for a registered 64bits source (from config file)
	r, found := registeredZBSources[source64]
	if !found {
		source64, err = xbee.Fix64address(source16)
		if err != nil {
			return fmt.Errorf("failed to fix the corrupted/missing 64-bits source address: %w", err)
		}
		log.Debugf("64-bits source address had been fixed: %s -> %s", source16, source64)
	} else if source16 != "" && source16 != "FFFF" && source16 != "FFEE" && source16 != "0000" {
		xbee.Record16bitsAddress(source64, source16)
	}

	u, err := url.Parse(viper.GetString("zb_redirection_url"))
	if err != nil {
		return errors.New("Can't parse redirection URL: " + err.Error())
	}

	q := u.Query()
	q.Add("room", r.Room)
	q.Add("type", r.Format)
	u.RawQuery = q.Encode()

	if resp, err := http.Post(u.String(), "text/plain", bytes.NewReader(data)); err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		var msg string
		if err != nil {
			msg = err.Error()
		} else {
			msg = resp.Status
		}

		return errors.New("Failed to redirect data to destination: " + msg)
	}

	return nil
}
