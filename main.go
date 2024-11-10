package main

import (
	"context"
	"trading-bot/core"
	"trading-bot/vendors/capital"
	"trading-bot/vendors/twelvedata"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core.ProcessQuotes()

	vendors := []core.QuoteStreamer{
		twelvedata.New(ctx),
		capital.New(ctx),
	}

	for _, vendor := range vendors {
		core.Inject(vendor)
	}

	core.HandleErrorsAndShutdown(cancel)
}
