package main

import (
	"fmt"
	"log"
	"trading-bot/core"

	"trading-bot/vendors/capital"
	"trading-bot/vendors/twelvedata"
	// "trading-bot/platform/config"
)

func main(){
    fmt.Println("Applaza trading bot with different API for creating and opening orders and fetching data. It integrates with different data providers and platform")

	//@TODO handle graceful shutdown of the listens when error happens
	quoteStreamers := []core.QuoteStreamer{twelvedata.QuoteStreamer{}, capital.QuoteStreamer{}}
	for _, streamer := range quoteStreamers {
		if err := streamer.StreamQuotes(); err != nil {
			log.Fatal(err)
		}
	}

	select {}
}