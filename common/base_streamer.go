package common

import (
	"context"
	"fmt"
	"log"
)

type BaseQuoteStreamer struct {
	Name       string
	QuotesChan chan Quote
	Ctx        context.Context
}

func (s *BaseQuoteStreamer) GetName() string {
	return s.Name
}

func (s *BaseQuoteStreamer) SetQuotesChannel(c chan Quote) {
	s.QuotesChan = c
}

// StreamQuotes is a placeholder method to be overridden by specific streamers.
func (s *BaseQuoteStreamer) StreamQuotes() error {
	return fmt.Errorf("[%s] StreamQuotes method not implemented", s.Name)
}

//open position
func (s *BaseQuoteStreamer) OpenPosition(Quote) error {
	return fmt.Errorf("[%s] openPosition method not implemented", s.Name)
}

func (s *BaseQuoteStreamer) PublishQuotes(quote Quote) {


	select {
	case s.QuotesChan <- quote:
		// log.Printf("[%s] Published quote: %v", s.Name, quote)
	case <-s.Ctx.Done():
		log.Printf("[%s] Context canceled, stopping publishing", s.Name)
	default:
		len := len(s.QuotesChan)
		log.Printf("[%s] Quote channel is full %v, dropping quote: %v", s.Name, len, quote)
	}
}

func (s *BaseQuoteStreamer) Log(format string, template ...interface{}) {
	log.Printf("["+s.Name+"] "+format, template...)
}
