package main

import (
	"trading-bot/core"
	"trading-bot/vendors/capital"
	"trading-bot/vendors/twelvedata"
)


func main() {

	vendors := []core.QuoteStreamer{twelvedata.New(), capital.New()}
	for _, vendor := range vendors {
		core.Inject(vendor)
	}

    core.HandleErrorsAndShutdown()
}
