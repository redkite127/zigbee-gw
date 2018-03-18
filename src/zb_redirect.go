package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type ZBSource struct {
	Room   string
	Format string
}

var registeredZBSources map[string]ZBSource

func ZBredirect(source string, data []byte) error {
	r, found := registeredZBSources[source]
	if !found {
		return errors.New("Received packet from unregistered device")
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
