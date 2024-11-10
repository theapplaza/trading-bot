package models

type PriceQuote struct{
	Price float64
	Symbol Symbol
}

type Symbol struct {
	Name string
	Ticker string
}
