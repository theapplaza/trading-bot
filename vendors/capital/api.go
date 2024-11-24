package capital

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trading-bot/common"
)

type CapitalApi struct {}

func NewCapitalApi() *CapitalApi {
	return &CapitalApi{}
}

func (api *CapitalApi) getPriceHistory(producer string, symbol string, count int, timeframe common.Timeframe) ([]common.Quote, error) {
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/prices/%s?resolution=%s&max=%d", activeConfig.ApiBaseUrl, symbol, string(timeframe), count+1), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cst", activeSession.cst)
	req.Header.Set("x-security-token", activeSession.securitytoken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get price history %s %s", err, err)
	}

	var result struct {
		Prices []struct {
			SnapshotTime string `json:"snapshotTime"`
			OpenPrice    struct {
				Bid float64 `json:"bid"`
				Ask float64 `json:"ask"`
			} `json:"openPrice"`
			ClosePrice struct {
				Bid float64 `json:"bid"`
				Ask float64 `json:"ask"`
			} `json:"closePrice"`
			HighPrice struct {
				Bid float64 `json:"bid"`
				Ask float64 `json:"ask"`
			} `json:"highPrice"`
			LowPrice struct {
				Bid float64 `json:"bid"`
				Ask float64 `json:"ask"`
			} `json:"lowPrice"`
		} `json:"prices"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var quotes []common.Quote
	for _, price := range result.Prices {
		timestamp, err := time.Parse("2006-01-02T15:04:05", price.SnapshotTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing timestamp: %v", err)
		}

		bidQuote := common.OhlcPriceQuote{
			Producer:   producer,
			Timeframe:  timeframe,
			QuoteType:  common.PriceQuoteTypeBid,
			HighPrice:  price.HighPrice.Bid,
			LowPrice:   price.LowPrice.Bid,
			OpenPrice:  price.OpenPrice.Bid,
			ClosePrice: price.ClosePrice.Bid,
			Symbol: common.Symbol{
				Name:   symbol,
				Ticker: symbol,
			},
			Timestamp: float64(timestamp.Unix()),
		}

		askQuote := common.OhlcPriceQuote{
			Producer:   producer,
			Timeframe:  timeframe,
			QuoteType:  common.PriceQuoteTypeAsk,
			HighPrice:  price.HighPrice.Ask,
			LowPrice:   price.LowPrice.Ask,
			OpenPrice:  price.OpenPrice.Ask,
			ClosePrice: price.ClosePrice.Ask,
			Symbol: common.Symbol{
				Name:   symbol,
				Ticker: symbol,
			},
			Timestamp: float64(timestamp.Unix()),
		}

		quotes = append(quotes, bidQuote, askQuote)

	}

	return quotes, nil
}

//openPosition
func (api *CapitalApi) openPosition(symbol string, quantity float64, side common.OrderSide) error {

	
	reqBody, _ := json.Marshal(map[string]interface{}{
		"epic":     symbol,
		"size": quantity,
		"direction": string(side),
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/positions", activeConfig.ApiBaseUrl), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cst", activeSession.cst)
	req.Header.Set("x-security-token", activeSession.securitytoken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot open position %s %s", err, err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error opening position: %v", response.Status)
	}

	return nil
}

//listPositions
func (api *CapitalApi) listPositions(s common.Symbol) ([]common.Position, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/positions", activeConfig.ApiBaseUrl), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cst", activeSession.cst)
	req.Header.Set("x-security-token", activeSession.securitytoken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot list positions %s %s", err, err)
	}

	var result struct {
		Positions []struct {
			Market struct {
				Epic string `json:"epic"`
			} `json:"market"`
			Position struct {
				Size float64 `json:"size"`
				Direction string `json:"direction"`
			} `json:"position"`
		} `json:"positions"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var positions []common.Position
	for _, position := range result.Positions {

		if position.Market.Epic != s.Name {
			continue
		}

		positions = append(positions, common.Position{
			Symbol:    common.Symbol{
				Name:   position.Market.Epic,
				Ticker: position.Market.Epic,
			},
			Quantity:  position.Position.Size,
			Side:      common.OrderSide(position.Position.Direction),
		})
	}

	return positions, nil
}