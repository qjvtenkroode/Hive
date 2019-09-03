package main

import (
	"log"
	"net/http"
)

func main() {
	store := new(InMemoryStore)
	seed := make(map[string]Asset)
	store.Assets = seed
	server := NewServer(store)

	if err := http.ListenAndServe(":80", server); err != nil {
		log.Fatalf("could not listen on port 80 %v", err)
	}
}
