package indicators

import (
	"fmt"

	"github.com/markcheno/go-talib"
)

/**
 * Implement this interface so that the core system can retrieve your implementation to work with for RSI indicator
 */

 var period = 14
 type Rsi interface {
	Get() float64
	Compute() float64
 }

func Compute(prices []float64) (result float64) {

	if len(prices) <= period {
		fmt.Println("Not enough data to calculate RSI")
		return 
	}

	values := talib.Rsi(prices, period)

	result = values[len(values)-1];

	fmt.Println("Rsi ", result)

	return
}
