package core

import (
	"log"
	"sync"
	"trading-bot/common"
)

var (
    streamerGroup sync.WaitGroup
    quoteChannel  = make(chan common.PriceQuote, 100) // Buffered channel to avoid blocking

)

func Inject(streamer common.QuoteStreamer) {
	streamer.SetQuotesChannel(quoteChannel)
	streamerGroup.Add(1)
	go func() {
		defer streamerGroup.Done()
		if err := streamer.StreamQuotes(); err != nil {
			log.Printf("[%s] Go routine returned with an error, so it has stopped", streamer.GetName())
			log.Printf("[%s] Error received - %v", streamer.GetName(), err)
		}
	}()

}

func ProcessQuotes() {
    go func() {
        for quote := range quoteChannel {
            log.Printf("Received quote: %v", quote)
            // Process the quote as needed
        }
    }()
}
