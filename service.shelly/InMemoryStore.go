package main

import (
	"errors"
)

// InMemoryStore is a small in-memory implementation of a Store
type InMemoryStore struct {
	States map[string]ShellyState
}

func (i *InMemoryStore) getState(id string) (ShellyState, error) {
	var err error
	asset, ok := i.States[id]
	if !ok {
		err = errors.New("asset not found in store")
	}
	return asset, err
}

func (i *InMemoryStore) storeState(id string, value ShellyState) error {
	var err error
	i.States[id] = value
	return err
}
