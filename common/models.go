package common

type Quote interface {
	IsQuote()
	GetSymbol() Symbol
	GetPrice() float64
	GetQuoteType() PriceQuoteType
	GetProducer() string
}

type Api interface {
	GetPriceHistory(producer string, symbol string, count int, timeframe Timeframe) ([]Quote, error)
	CreateOrder(symbol string, price float64, quantity float64, side OrderSide) error
}

type OhlcPriceQuote struct {
	Producer   string
	Timeframe     Timeframe
	QuoteType  PriceQuoteType
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

type Position struct {
	Symbol   Symbol
	Quantity float64
	Side     OrderSide
}

func (p OhlcPriceQuote) IsQuote() {}
func (p PriceQuote) IsQuote()       {}
func (p OhlcPriceQuote) GetSymbol() Symbol {
	return p.Symbol
}
func (p PriceQuote) GetSymbol() Symbol {
	return p.Symbol
}
func (p OhlcPriceQuote) GetPrice() float64 {
	return p.ClosePrice
}
func (p PriceQuote) GetPrice() float64 {
	return p.Price
}
func (p OhlcPriceQuote) GetQuoteType() PriceQuoteType {
	return p.QuoteType
}

//@NOTE: how do we handle realtime quotes?
func (p PriceQuote) GetQuoteType() PriceQuoteType {
	return PriceQuoteTypeBid
}

//get producer
func (p OhlcPriceQuote) GetProducer() string {
	return p.Producer
}

func (p PriceQuote) GetProducer() string {
	return p.Producer
}

// QuoteStreamer defines the interface for streaming quotes.
type QuoteStreamer interface {
	SetQuotesChannel(chan Quote)
	GetName() string

	//Client must implement as this is the try point to the vendor
	StreamQuotes() error
	GetStrategies() []SignalStrategy
	OpenPosition(Quote) error
}
