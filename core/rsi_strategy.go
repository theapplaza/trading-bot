package core

import (
	"log"
	"trading-bot/common"

	"github.com/markcheno/go-talib"
)

type RsiStrategy struct {
	Name   common.SignalStrategy
	Period int
	BidSignalLevel  float64
	AskSignalLevel  float64
}

func NewRsiStrategy(period int, bidLevel float64, askLevel float64) *RsiStrategy {
	return &RsiStrategy{
		Name:   common.Rsi,
		Period: period,
		BidSignalLevel: bidLevel,
		AskSignalLevel: askLevel,
	}
}

func (rsi *RsiStrategy) Check(currentQuote common.Quote) (current float64, ok bool) {

	data := dataStore.GetData(currentQuote.GetProducer(), currentQuote.GetSymbol())
	if data == nil {
		return 0, false
	}

	if len(data) < (rsi.Period+1) * 2 {
		// log.Printf("Not enough data for RSI calculation for %s %s", currentQuote.GetProducer(), currentQuote.GetSymbol().Name)
		return 0, false
	}

	var prices []float64
	for i := len(data) - ((rsi.Period+1)*2); i < len(data); i++ {
		if currentQuote.GetQuoteType() == data[i].GetQuoteType() {
			prices = append(prices, data[i].GetPrice())
		}
	}

	if len(prices) < rsi.Period {
		return 0, false
	}

	rsiValues := talib.Rsi(prices, rsi.Period)
    current = rsiValues[len(rsiValues)-1]

	if currentQuote.GetQuoteType() == "bid" {
		ok = current > rsi.BidSignalLevel
	} else {
		ok = current < rsi.AskSignalLevel 
	}

	log.Printf("RSI for %s and kind %s: %.2f", currentQuote.GetSymbol().Name, currentQuote.GetQuoteType(), current)
	
	return current, ok
}