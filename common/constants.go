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

type SymbolClass string
const (
	Currency SymbolClass = "currency"
	Stock SymbolClass = "stock"
	Crypto SymbolClass = "crypto"
)

type Timeframe string
const (
	OneMinute Timeframe = "MINUTE"
	FiveMinutes Timeframe = "MINUTE_5" 
	Hour Timeframe = "HOUR"
	Day Timeframe = "DAY"
	Week Timeframe = "WEEK"
	Month Timeframe = "MONTH"
)

type OrderSide string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"
)

type OrderType string
const (
	// Market OrderType = "MARKET"
	Limit OrderType = "LIMIT"
	Stop OrderType = "STOP"
	// StopLimit OrderType = "STOP_LIMIT"
)