package main

import (
	"log"
	"net/http"
)

func main() {
	store := new(InMemoryStore)
	store.Assets = make(map[string]Asset)
	server := NewServer(store)

	if err := http.ListenAndServe(":80", server); err != nil {
		log.Fatalf("could not listen on port 80 %v", err)
	}
}
