package core

import (
	"log"
	"sync"
	"trading-bot/common"
)

type VendorData struct {
	Quotes []common.Quote
}

type DataStore struct {
	mu  sync.Mutex
	storage map[string]*VendorData
}

func NewDataStore() *DataStore {
	return &DataStore{
		storage: make(map[string]*VendorData),
	}
}

func (ds *DataStore) AddVendor(vendorName string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.storage[vendorName] = &VendorData{
		Quotes: []common.Quote{},
	}
}
func (ds *DataStore) AddData(vendorName string, data common.Quote) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		v.Quotes = append(v.Quotes, data)
		log.Printf("added realtime data for %s data is now %v", vendorName, v.Quotes)

	}
}

func (ds *DataStore) GetData(vendorName string) []common.Quote {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		return v.Quotes
	}
	return nil
}
