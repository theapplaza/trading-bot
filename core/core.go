package core

import (
	"log"
	"reflect"
	"sync"
	"trading-bot/common"
)

var (
	streamerGroup sync.WaitGroup
	quoteChannel  = make(chan interface{}, 100) // Buffered channel to avoid blocking

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
			quoteType := reflect.TypeOf(quote)
			q := reflect.ValueOf(quote)

			switch quoteType {
			case reflect.TypeOf(common.PriceQuote{}):
				// Process PriceQuote
				log.Println("Processing PriceQuote")
				log.Printf("Received from %v, instrument: %v, value=%f", q.FieldByName("Producer"), q.FieldByName("Symbol").FieldByName("Name"), q.FieldByName("Price").Float())
			case reflect.TypeOf(common.PeriodPriceQuote{}):
				// Process ForexQuote
				log.Println("Processing PeriodPriceQuote")
				log.Printf("Received from %v, instrument: %v, pricetype: %s, value=%f", q.FieldByName("Producer"), q.FieldByName("Symbol").FieldByName("Name"), q.FieldByName("QuoteType"), q.FieldByName("ClosePrice").Float())

			default:
				log.Printf("Unknown quote type: %v", quoteType)
			}
			// log.Printf(
			// 	"Received from %v, instrument: %v, value=%f",
			// 	quote.Producer,
			// 	quote.Symbol.Name,
			// 	quote.Price,
			// )
			// Process the quote as needed
		}
	}()
}
