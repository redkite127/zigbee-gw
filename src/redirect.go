package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

var registeredRedirections map[string]interface{}

func redirect(source string, data []byte) error {
	uStr, found := registeredRedirections[source].(string)
	if !found {
		return errors.New("Received packet from unregistered device")
	}

	u, err := url.Parse(uStr)
	if err != nil {
		return errors.New("Can't parse redirection URL: " + err.Error())
	}

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
