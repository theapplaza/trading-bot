package core

import (
	"log"
	"sync"
	"trading-bot/common"
)

type VendorData struct {
	RealtimeData []common.PriceQuote
	PeriodPriceData []common.PeriodPriceQuote //@TODO: we are assumming we are working with a single period for now
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
		RealtimeData: make([]common.PriceQuote, 0),
		PeriodPriceData: make([]common.PeriodPriceQuote, 0),
	}
}

func (ds *DataStore) AddRealtimeData(vendorName string, data common.PriceQuote) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		v.RealtimeData = append(v.RealtimeData, data)
		log.Printf("added realtime data for %s data is now %v", vendorName, v.RealtimeData)
	}
}

func (ds *DataStore) GetRealtimeData(vendorName string) []common.PriceQuote {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		return v.RealtimeData
	}
	return nil
}

func (ds *DataStore) AddPeriodPriceData(vendorName string, data common.PeriodPriceQuote) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		v.PeriodPriceData = append(v.PeriodPriceData, data)
	}
}

func (ds *DataStore) GetPeriodPriceData(vendorName string) []common.PeriodPriceQuote {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		return v.PeriodPriceData
	}
	return nil
}