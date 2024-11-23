package core

import (
	"log"
	"sync"
	"trading-bot/common"
)

var (
	streamerGroup sync.WaitGroup
	quoteChannel  = make(chan common.Quote, 100) // Buffered channel to avoid blocking
	dataStore     = NewDataStore()
	strategies    = make(map[common.SignalStrategy]Strategy)
)

func init() {
	strategies[common.Rsi] = NewRsiStrategy(14, 70, 30)
}

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
				dataStore.AddPriceQuote(q.Producer, q)
			case common.PeriodPriceQuote:
				dataStore.AddPeriodPriceQuote(q.Producer, q)
			default:
				log.Printf("Unknown quote type: %T", quote)
			}

			if strategy, ok := strategies[common.Rsi]; ok {
				 strategy.Check(quote); 
			}
		}
	}()
}
