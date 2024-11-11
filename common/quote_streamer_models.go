package common

type PeriodPriceQuote struct {
	Producer   string
	Period     string
	QuoteType  string
	HighPrice  float64
	LowPrice   float64
	OpenPrice  float64
	ClosePrice float64
	Symbol     Symbol
	Timestamp  float64
}

type PriceQuote struct {
	Producer  string
	Price     float64
	Symbol    Symbol
	Timestamp float64
}

type Symbol struct {
	Name   string
	Ticker string
}

// QuoteStreamer defines the interface for streaming quotes.
type QuoteStreamer interface {
	SetQuotesChannel(chan interface{})
	GetName() string

	//Client must implement as this is the try point to the vendor
	StreamQuotes() error
}
