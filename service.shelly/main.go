package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/eclipse/paho.mqtt.golang"
)

func main() {
	connOpts := mqtt.NewClientOptions()
	connOpts.AddBroker("tcp://mqtt.qkroode.nl:1883")
	connOpts.SetClientID("hive.shelly")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername("hivemind")
	connOpts.SetPassword(",33TnJLPMy>VNax")
	connOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert})

	client := mqtt.NewClient(connOpts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalf("mqtt failed: %v", token.Error())
	}

	store := new(InMemoryStore)
	store.States = make(map[string]ShellyState)

	server := NewServer(store, client, token)
	if server.token = client.Subscribe("shellies/+/relay/0", byte(0), server.onMessageReceived); server.token.Wait() && server.token.Error() != nil {
		log.Fatalf("mqtt failed: %v", server.token.Error())
	}

	if err := http.ListenAndServe(":80", server); err != nil {
		log.Fatalf("could not listen on port 80 %v", err)
	}
}
