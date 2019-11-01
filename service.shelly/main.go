package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/eclipse/paho.mqtt.golang"
)

func main() {
	port := flag.String("p", "80", "port to use when running")
	mqttUser := flag.String("u", "", "mqtt username")
	mqttPassword := flag.String("pw", "", "mqtt password")
	mqttHost := flag.String("h", "", "mqtt host and port <host>:<port>")
	flag.Parse()

	connOpts := mqtt.NewClientOptions()
	connOpts.AddBroker("tcp://" + *mqttHost)
	connOpts.SetClientID("hive.shelly")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername(*mqttUser)
	connOpts.SetPassword(*mqttPassword)
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

	if err := http.ListenAndServe(":"+*port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", *port, err)
	}
}
