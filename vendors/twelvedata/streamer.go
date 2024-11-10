package twelvedata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"trading-bot/common"

	"github.com/gorilla/websocket"
)

type TwelveDataStreamer struct {
	common.BaseQuoteStreamer
	url string
}

func New(ctx context.Context) TwelveDataStreamer {
	_url := fmt.Sprintf("wss://%v/quotes/price?apikey=%v", activeConfig.DataStreamUrl, activeConfig.ApiKey)

	return TwelveDataStreamer{
		BaseQuoteStreamer: common.BaseQuoteStreamer{
			Name: "TwelveData",
			Ctx: ctx,
		},
		url: _url,
	}
}


func (s TwelveDataStreamer) StreamQuotes() (err error) {

	con, _, err := websocket.DefaultDialer.Dial(s.url, nil)
	if err != nil {
		return err
	}
	defer con.Close()

	subscribePayload := map[string]interface{}{
		"action": "subscribe",
		"params": map[string]interface{}{
			"symbols": activeConfig.Instruments,
		},
	}
	con.WriteJSON(subscribePayload)

	return s.listen(con)
}

func (s TwelveDataStreamer) listen(con *websocket.Conn) (err error) {

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
				log.Println("error unmarshalling response:", err)
				continue
			}

			// Handle the response based on the destination field
			switch response["event"] {
			case "connection":
				return err
			case "subscribe-status":
				if err := s.handleSubscriptionResponse(response); err != nil {
					return err
				}
			case "price":
				s.handleQuoteUpdateResponse(response)
			default:
				return fmt.Errorf("unhandled message: %v", response)
			}
		}
	}
}

func (s TwelveDataStreamer) handleSubscriptionResponse(response map[string]interface{}) error {
	status := response["status"].(string)
	if status != "ok" {
		return fmt.Errorf("subscription failed: %v", response)
	}
	return nil
}

func (s TwelveDataStreamer) handleQuoteUpdateResponse(payload map[string]interface{}) {
	epic := payload["symbol"].(string)
	price := payload["price"].(float64)
	timestamp := payload["timestamp"].(float64)

	quote := common.PriceQuote{
		Price: price,
		Symbol: common.Symbol{
			Name:   epic,
			Ticker: epic,
		},
	}
	s.PublishQuotes(quote)

	log.Printf("TwelveData: Price update for %s - Price: %f, Timestamp: %d", epic, price, int64(timestamp))
}
