package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ZBSource struct {
	Room   string
	Format string
}

// Recording last 16bits address in case a 64bit address gets corrupted... (it happend)
var source16bits = map[string]string{}
var registeredZBSources map[string]ZBSource

func ZBredirect(source64, source16 string, data []byte) error {
	// looking for a registered 64bits source (from config file)
	r, found := registeredZBSources[source64]
	if !found {
		// if not found, the source address may have been corrupted, looking inside previous
		// successfully matched 64bits sources for a 16bits source match. (the 16bits address is given when joinging a newtwork)
		fixed64, found2 := source16bits[source16]
		if !found2 {
			return errors.New("Received packet from unregistered device")
		}
		// no need to check found or not, it's a 64bits source comming from the config file!
		r = registeredZBSources[fixed64]
		log.Infof("64bits sources was corrupted, fixed: %s -> %s\n", source64, source16)
	} else {
		// Got an immediate match, recording its 16bits address
		source16bits[source16] = source64
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
