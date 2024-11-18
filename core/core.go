package core

import (
	"log"
	"sync"
	"trading-bot/common"
)

var (
	streamerGroup sync.WaitGroup
	quoteChannel  = make(chan interface{}, 100) // Buffered channel to avoid blocking
	dataStore	 = NewDataStore()
)

func Inject(streamer common.QuoteStreamer) {
	vendorName := streamer.GetName()
	dataStore.AddVendor(vendorName)
	streamer.SetQuotesChannel(quoteChannel)
	streamerGroup.Add(1)
	go func() {
		defer streamerGroup.Done()
		if err := streamer.StreamQuotes(); err != nil {
			log.Printf("[%s] Go routine stopped with error: %v", vendorName, err)
		}
	}()
}

func ProcessQuotes() {
	go func() {
		for quote := range quoteChannel {
			switch q := quote.(type) {
			case common.PriceQuote:
				dataStore.AddRealtimeData(q.Producer, q)
			case common.PeriodPriceQuote:
				dataStore.AddPeriodPriceData(q.Producer, q)
			default:
				log.Printf("Unknown quote type: %T", quote)
			}
		}
	}()
}
