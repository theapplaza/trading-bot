package core

import (
	"sync"
	"trading-bot/common"
)

type VendorData struct {
	Quotes map[string][]common.Quote
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
		Quotes: make(map[string][]common.Quote),
	}
}

//AddPriceQuote adds a price quote to the data store
func (ds *DataStore) AddPriceQuote(vendorName string, data common.PriceQuote) {
	ds.AddData(vendorName, data.Symbol, data)
}

//AddOhlcPriceQuote adds a period price quote to the data store
func (ds *DataStore) AddOhlcPriceQuote(vendorName string, data common.OhlcPriceQuote) {
	ds.AddData(vendorName, data.Symbol, data)
}

func (ds *DataStore) AddData(vendorName string, symbol common.Symbol, data common.Quote) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		v.Quotes[symbol.Name] = append(v.Quotes[symbol.Name], data)
	}
}

//Gets data fora single symbol
func (ds *DataStore) GetData(vendorName string, symbol common.Symbol) []common.Quote {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		return v.Quotes[symbol.Name]
	}
	return nil
}

func (ds *DataStore) GetSymbols(vendorName string) []common.Symbol {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if v, ok := ds.storage[vendorName]; ok {
		var symbols []common.Symbol
		for symbol := range v.Quotes {
			symbols = append(symbols, common.Symbol{Name: symbol})
		}
		return symbols
	}
	return nil
}