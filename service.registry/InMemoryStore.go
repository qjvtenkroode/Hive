package main

import (
	"errors"
	"sort"
)

// InMemoryStore is a small in-memory implementation of a Store
type InMemoryStore struct {
	Assets map[string]Asset
}

func (i *InMemoryStore) getAllAssets() []Asset {
	var assets []Asset
	for _, asset := range i.Assets {
		assets = append(assets, asset)
	}
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Identifier < assets[j].Identifier
	})
	return assets
}

func (i *InMemoryStore) getAsset(id string) (Asset, error) {
	var err error
	asset, ok := i.Assets[id]
	if !ok {
		err = errors.New("asset not found in store")
	}
	return asset, err
}

func (i *InMemoryStore) storeAsset(id string, value Asset) error {
	var err error
	i.Assets[id] = value
	return err
}

func (i *InMemoryStore) deleteAsset(id string) error {
	var err error
	if _, ok := i.Assets[id]; ok {
		delete(i.Assets, id)
	} else {
		err = errors.New("asset not found in store")
	}
	return err
}
