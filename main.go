package main

import (
	"fmt"
	"log"
	"reflect"
	"trading-bot/core"
	"trading-bot/vendors/capital"
	"trading-bot/vendors/twelvedata"
)

func main() {
	fmt.Println("Applaza trading bot integrates with different data providers and platform to monitor and execute orders")

	errorChannel := core.GetErrorChannel()
	quoteStreamers := []core.QuoteStreamer{
		twelvedata.New(errorChannel),
		capital.New(errorChannel),
	}
	for _, streamer := range quoteStreamers {
		if err := streamer.StreamQuotes(); err != nil{
			log.Printf("cannot start streaming quotes for %s", reflect.TypeOf(streamer).Name())
		}
	}

	core.HandleErrors()
}
