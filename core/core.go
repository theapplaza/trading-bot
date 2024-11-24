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
	streamers	 = make(map[string]common.QuoteStreamer)
)

func init() {
	strategies[common.Rsi] = NewRsiStrategy(14, 70, 30)
}

func Inject(streamer common.QuoteStreamer) {
	vendorName := streamer.GetName()
	dataStore.AddVendor(vendorName)
	streamer.SetQuotesChannel(quoteChannel)
	streamers[vendorName] = streamer
	streamerGroup.Add(1)
	go func() {
		defer streamerGroup.Done()
		//todo remove streamer from group
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
			case common.OhlcPriceQuote:
				dataStore.AddOhlcPriceQuote(q.Producer, q)
			default:
				log.Printf("Unknown quote type: %T", quote)
			}

			if strategy, ok := strategies[common.Rsi]; ok {
				if _, ok := strategy.Check(quote); ok  {
					//get the streamer and create the order
					streamer := GetStreamer(quote)
					if streamer != nil {
						log.Println("Opening position")
						if err := streamer.OpenPosition(quote); err != nil {
							log.Printf("[%s] Error opening position: %v", streamer.GetName(), err)
						}else{
							log.Printf("[%s] Position opened", streamer.GetName())
						}
					}
				}
			}
		}
	}()
}

func GetStreamer(q common.Quote) common.QuoteStreamer {
	return streamers[q.GetProducer()]
}
