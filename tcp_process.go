package main

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func processTCPrequests(stop <-chan bool, stopped chan<- bool) {
	server := &http.Server{
		Addr: ":" + viper.GetString("tcp.port"),
		// Handler: http.TimeoutHandler(nil, 30*time.Second, "503 Service Unavailable"),
	}

	http.HandleFunc("/zigbee", requestHandler)

	go func() {
		log.Infof("waiting TCP requests: 0.0.0.0:%d", viper.GetInt("tcp.port"))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Infof("interrupted... no more TCP request to process")
	stopped <- true
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	destination := v.Get("destination")
	room := v.Get("room")
	if destination == "" && room != "" {
		for a, s := range registeredZBSources {
			if s.Room == room {
				destination = a
			}
		}
	}

	if destination == "" {
		log.Infoln("received TCP request for unregistered device")
		http.NotFound(w, r)
		return
	}

	data := v.Get("data")

	if err := TCPredirect(destination, data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
