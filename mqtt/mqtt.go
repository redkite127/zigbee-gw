package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var Client mqtt.Client

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.WithField("backend", "mqtt").Infoln("connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.WithField("backend", "mqtt").Warningln("connection lost: %w", err)
	//TODO quit gracefully?
}

func InitMQTT(host string, port int, username, password string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetClientID("home-hub")
	opts.SetUsername(username)
	opts.SetPassword(password)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Publish(topic string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[MQTT] failed to encode the payload to json: %w", err)
	}

	token := Client.Publish(topic, 0, false, jsonData)
	if !token.WaitTimeout(1*time.Second) || token.Error() != nil {
		return fmt.Errorf("[MQTT] failed to publish the payload: %w", token.Error())
	}

	return nil
}
