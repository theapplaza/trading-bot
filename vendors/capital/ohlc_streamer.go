package capital

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"trading-bot/common"

	"github.com/gorilla/websocket"
)

type OhlcStreamer struct {
	common.BaseQuoteStreamer
	period string
}

func NewOhlcStreamer(ctx context.Context) *OhlcStreamer {
	return &OhlcStreamer{
		BaseQuoteStreamer: common.BaseQuoteStreamer{
			Name: "Capital OHLC",
			Ctx:  ctx,
		},
		period: "MINUTE",
	}
}

func (s *OhlcStreamer) GetStrategies() []common.SignalStrategy {
	return nil
}

func (s *OhlcStreamer) StreamQuotes() (err error) {

	//ensure that authentication is done
	if err = authenticate(); err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", activeSession.StreamingHost, "connect")
	con, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer con.Close()

	//subcribe to price update
	request := map[string]interface{}{
		"destination":   "OHLCMarketData.subscribe",
		"correlationId": rand.Int(),
		"cst":           activeSession.cst,
		"securityToken": activeSession.securitytoken,
		"payload": map[string]interface{}{
			"epics":       strings.Split(activeConfig.Instruments, ","),
			"resolutions": []string{s.period},
			"type":        "classic",
		},
	}
	err = con.WriteJSON(request)

	return s.listen(con)
}

func (s *OhlcStreamer) listen(con *websocket.Conn) (err error) {

	defer con.Close()

	for {

		select {
		case <-s.Ctx.Done():
			return s.Ctx.Err()
		default:
			_, message, err := con.ReadMessage()
			if err != nil {
				return err
			}

			var response map[string]interface{}
			if err := json.Unmarshal(message, &response); err != nil {
				log.Println("Error unmarshalling response:", err)
				continue
			}

			// Handle the response based on the destination field
			switch response["destination"] {
			case "OHLCMarketData.subscribe":
				if err := s.handleSubscriptionResponse(response); err != nil {
					return err
				}
			case "ohlc.event":
				s.handleQuoteUpdateResponse(response)
			default:
				return fmt.Errorf("unhandled message: %v", response)
			}
		}
	}
}

func (s *OhlcStreamer) handleSubscriptionResponse(response map[string]interface{}) error {
	status := response["status"].(string)
	if status != "OK" {
		return fmt.Errorf("subscription error: %v", response)
	}

	payload := response["payload"].(map[string]interface{})
	subcriptions := payload["subscriptions"].(map[string]interface{})

	for _, value := range subcriptions {
		v := value.(string)
		if strings.Contains(v, "ERROR") {
			return fmt.Errorf("subscription error: %v", v)
		}
	}

	return nil
}

func (s *OhlcStreamer) handleQuoteUpdateResponse(response map[string]interface{}) {
	payload := response["payload"].(map[string]interface{})

	quote := common.PeriodPriceQuote{
		Producer:   s.GetName(),
		Period:     s.period,
		QuoteType:  payload["priceType"].(string),
		HighPrice:  payload["h"].(float64),
		LowPrice:   payload["l"].(float64),
		OpenPrice:  payload["o"].(float64),
		ClosePrice: payload["c"].(float64),
		Symbol: common.Symbol{
			Name:   payload["epic"].(string),
			Ticker: payload["epic"].(string),
		},
		Timestamp: payload["t"].(float64),
	}

	s.PublishQuotes(quote)

	// log.Printf("Capital HOLC: Price update for %s - PriceType: %s, Open: %f, Close: %.2f, Timestamp: %f", quote.Symbol.Name, quote.QuoteType, quote.OpenPrice, quote.ClosePrice, float64(quote.Timestamp))
}
