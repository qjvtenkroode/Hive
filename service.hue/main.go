package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "80", "port to use when running")
	bridgeAddress := flag.String("ha", "", "hue bridge address")
	bridgeToken := flag.String("ht", "", "hue bridge token")
	readme := flag.String("readme", "README.md", "custom readme file")
	flag.Parse()

	bridge := new(HueBridge)
	bridge.Token = *bridgeToken
	bridge.Address = *bridgeAddress

	store := new(InMemoryStore)
	store.States = make(map[string]HueLightState)

	server := NewServer(store, *bridge, *readme)

	server.pollState()

	if err := http.ListenAndServe(":"+*port, server); err != nil {
		log.Fatalf("could not listen on port 80 %v", err)
	}
}
