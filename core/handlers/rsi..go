package handlers

import (
	"fmt"

	"github.com/markcheno/go-talib"
)

func Compute(prices []float64) (result []float64) {
	if len(prices) < 15 {
		fmt.Println("Not enough data to calculate RSI")
		return 
	}
	result = talib.Rsi(prices, 14)

	fmt.Println("Rsi ", result)
	return
}
