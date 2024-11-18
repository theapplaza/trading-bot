package core

import (
	"trading-bot/common"

	"github.com/markcheno/go-talib"
)

type RsiStrategy struct {
	Name   common.SignalStrategy
	Period int
	Level  float64
}

func NewRsiStrategy(period int, level float64) *RsiStrategy {
	return &RsiStrategy{
		Name:   common.Rsi,
		Period: period,
		Level:  level,
	}
}

func (rsi *RsiStrategy) Pass(vendorName string, symbol common.Symbol) (current float64, ok bool) {

	//compute rsi for each symbol in store

		data := dataStore.GetData(vendorName, symbol)
		if data == nil {
			return 0, false
		}

		var closePrices []float64
		for _, quote := range data {
			if priceQuote, ok := quote.(common.PriceQuote); ok {
				closePrices = append(closePrices, priceQuote.Price)
			}
		}

		if len(closePrices) < rsi.Period+1 {
			return 0, false
		}

		rsiValues := talib.Rsi(closePrices, rsi.Period)
		current = rsiValues[len(rsiValues)-1]

		//log
		// log.Printf("RSI for %s: %s: %f", vendorName, symbol.Name, currentRsi)
		return current, current > rsi.Level

}
