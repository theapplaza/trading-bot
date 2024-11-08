package core

 type Rsi struct {
	Period Timeframe
	Value float64
	Prices []float64
 }

 type IRsi interface {
	Compute(period int) float64
 }