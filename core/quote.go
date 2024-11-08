package core

type QuoteImpl interface {
	Get() []float64
}

type Quote struct {
	Symbol Symbol
	Price float64
}