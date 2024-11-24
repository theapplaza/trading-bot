package common

type PriceQuoteType string
type SignalStrategy int

const (
	PriceQuoteTypeBid PriceQuoteType = "bid"
	PriceQuoteTypeAsk PriceQuoteType = "ask"
)

const (
	Rsi SignalStrategy = iota
	Macd
)