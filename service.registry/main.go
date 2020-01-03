package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "80", "port to use when running")
	readme := flag.String("readme", "README.md", "custom readme file")
	flag.Parse()

	store := new(InMemoryStore)
	store.Assets = make(map[string]Asset)
	server := NewServer(store, *readme)

	if err := http.ListenAndServe(":"+*port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", *port, err)
	}
}
