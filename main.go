package main

import (
	"fmt"
	"trading-bot/vendors/capital"
	// "trading-bot/core"
	// "trading-bot/platform/config"
)

func main(){
    fmt.Println("Applaza trading bot with different API for creating and opening orders and fetching data. It integrates with different data providers and platform")
	capital.Start()
}