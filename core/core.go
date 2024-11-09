package core

import (
	"log"
	"sync"
)

type QuoteStreamer interface {
	GetName() string
	StreamQuotes() error
}

var streamersGroup sync.WaitGroup

func Inject(streamer QuoteStreamer) {
	streamersGroup.Add(1)
	go func() {
		defer streamersGroup.Done()
		if err := streamer.StreamQuotes(); err != nil {
			log.Printf("[%s] Go routine returned with an error, so it has stopped", streamer.GetName())
			log.Printf("[%s] Error received - %v", streamer.GetName(), err)
		}
	}()
}
