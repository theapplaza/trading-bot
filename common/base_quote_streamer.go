package common

import (
	"context"
	"fmt"
	"log"
)

type BaseQuoteStreamer struct {
	Name       string
	QuotesChan chan PriceQuote
	Ctx        context.Context
}

func (s BaseQuoteStreamer) GetName() string {
	return s.Name
}

func (s BaseQuoteStreamer) SetQuotesChannel(c chan PriceQuote) {
	s.QuotesChan = c
}

// StreamQuotes is a placeholder method to be overridden by specific streamers.
func (s BaseQuoteStreamer) StreamQuotes() error {
	return fmt.Errorf("[%s] StreamQuotes method not implemented", s.Name)
}

func (s BaseQuoteStreamer) PublishQuotes(quote PriceQuote) {
	s.QuotesChan <- quote
	log.Println("publishing quotes")
	select {
    case s.QuotesChan <- quote:
        log.Printf("[%s] Published quote: %v", s.Name, quote)
    case <-s.Ctx.Done():
        log.Printf("[%s] Context canceled, stopping publishing", s.Name)
    default:
		len := len(s.QuotesChan)
        log.Printf("[%s] Quote channel is full %v, dropping quote: %v", s.Name,len, quote)
    }
}

func (s BaseQuoteStreamer) Log(format string, template ...interface{}) {
	log.Printf("["+s.Name+"] "+format, template...)
}
