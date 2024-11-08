package main

import (
	"trading-bot/core"
	"trading-bot/vendors/capital"
	"trading-bot/vendors/twelvedata"
)

func main() {
	// Start handling errors
	core.HandleErrors()

	twelvedata.New(core.GetErrorChannel()).StreamQuotes()
	capital.New(core.GetErrorChannel()).StreamQuotes()
	
	// Keep the main function running
	select {}
}
