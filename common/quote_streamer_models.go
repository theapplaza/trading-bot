package common

type PriceQuote struct{
	Price float64
	Symbol Symbol
}

type Symbol struct {
	Name string
	Ticker string
}

// QuoteStreamer defines the interface for streaming quotes.
type QuoteStreamer interface {
    SetQuotesChannel(chan PriceQuote)
	GetName() string

	 //Client must implement as this is the try point to the vendor
    StreamQuotes() error
}
