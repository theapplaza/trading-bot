package common

type Quote interface {
	IsQuote()
	GetSymbol() Symbol
	GetPrice() float64
	GetQuoteType() string
	GetProducer() string
}

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

func (p PeriodPriceQuote) IsQuote() {}
func (p PriceQuote) IsQuote()       {}
func (p PeriodPriceQuote) GetSymbol() Symbol {
	return p.Symbol
}
func (p PriceQuote) GetSymbol() Symbol {
	return p.Symbol
}
func (p PeriodPriceQuote) GetPrice() float64 {
	return p.ClosePrice
}
func (p PriceQuote) GetPrice() float64 {
	return p.Price
}
func (p PeriodPriceQuote) GetQuoteType() string {
	return p.QuoteType
}

//@NOTE: how do we handle realtime quotes?
func (p PriceQuote) GetQuoteType() string {
	return "bid"
}

//get producer
func (p PeriodPriceQuote) GetProducer() string {
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
}

type SignalStrategy int

const (
	Rsi SignalStrategy = iota
	Macd
)