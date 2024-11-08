package core

type QuoteStreamer interface {
	StreamQuotes() error
}