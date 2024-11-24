package capital

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trading-bot/common"
)

func getPriceHistory(producer string, symbol string, period int, resolution string) ([]common.Quote, error) {
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("%s/prices/%s?resolution=%s&max=%d", activeConfig.ApiBaseUrl, symbol, resolution, period+1), nil)
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
			Timeframe:     common.Timeframe(resolution),
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
			Timeframe:     common.Timeframe(resolution),
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
