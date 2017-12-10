package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

type Redirect struct {
	Room   string
	Format string
}

var registeredRedirections map[string]Redirect

func redirect(source string, data []byte) error {
	r, found := registeredRedirections[source]
	if !found {
		return errors.New("Received packet from unregistered device")
	}

	u, err := url.Parse("http://10.161.0.130:2001/sensors")
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
