package capital

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"trading-bot/common"

	"github.com/gorilla/websocket"
)

type OhlcStreamer struct {
	common.BaseQuoteStreamer
	common.Timeframe
}

var api *CapitalApi 

func NewOhlcStreamer(ctx context.Context) *OhlcStreamer {
	api = NewCapitalApi()
	return &OhlcStreamer{
		common.BaseQuoteStreamer{
			Name: "Capital OHLC",
			Ctx:  ctx,
		},
		common.OneMinute,
	}
}

func (s *OhlcStreamer) GetStrategies() []common.SignalStrategy {
	return nil
}

func (s *OhlcStreamer) StreamQuotes() (err error) {

	for {
		if err = s.connectAndStream(); err != nil {
			if errors.Is(err, ErrRecoverable) {
				s.Log("Recoverable error connecting to streamer:", err)
				time.Sleep(30 * time.Second)
			} else {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (s *OhlcStreamer) connectAndStream() (err error) {
	//ensure that authentication is done
	if err = authenticate(true); err != nil {
		return fmt.Errorf("%w: %v", ErrNonRecoverable, err)
	}

	//populate the base starting values
	s.SetHistoricalData()

	url := fmt.Sprintf("%s%s", activeSession.StreamingHost, "connect")
	con, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNonRecoverable, err)
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
			"resolutions": []string{string(s.Timeframe)},
			"type":        "classic",
		},
	}
	err = con.WriteJSON(request)

	return s.listen(con)
}

func (s *OhlcStreamer) listen(con *websocket.Conn) (err error) {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {

		select {
		case <-s.Ctx.Done():
			return s.Ctx.Err()
		case <-ticker.C:
			if err := con.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error sending ping:", err)
				return fmt.Errorf("%w: %v", ErrRecoverable, err)
			}
		default:
			_, message, err := con.ReadMessage()
			if err != nil {
				s.Log("Error reading message:", err)
				// if websocket.IsCloseError(err) {
				return fmt.Errorf("%w: %v", ErrRecoverable, err)
				// }
				// return fmt.Errorf("%w: %v", ErrNonRecoverable, err)
			}

			if err := s.processMessage(message); err != nil {
				s.Log("Error processing message:", err)
			}
		}
	}
}

// implement processMessage
func (s *OhlcStreamer) processMessage(message []byte) error {

	var response map[string]interface{}
	if err := json.Unmarshal(message, &response); err != nil {
		s.Log("Error unmarshalling response:", err)
		return nil
	}

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

	return nil
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

	s.Log("Subscription successful")

	return nil
}

func (s *OhlcStreamer) handleQuoteUpdateResponse(response map[string]interface{}) {
	payload := response["payload"].(map[string]interface{})

	quote := common.OhlcPriceQuote{
		Producer:   s.GetName(),
		Timeframe:  s.Timeframe,
		QuoteType:  common.PriceQuoteType(payload["priceType"].(string)),
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
}

func (s *OhlcStreamer) SetHistoricalData() error {

	instruments := strings.Split(activeConfig.Instruments, ",")

	for _, instrument := range instruments {
		symbolQuotes, _ := api.getPriceHistory(s.GetName(), instrument, 15, s.Timeframe)
		for _, quote := range symbolQuotes {
			s.PublishQuotes(quote)
		}
	}

	return nil
}

/**
 * @TODO: Abstract the open position logic to the common package like PublishQuotes
 */
func (s *OhlcStreamer) OpenPosition(q common.Quote) error {

	var side common.OrderSide
	if q.GetQuoteType() == common.PriceQuoteTypeBid {
		side = common.Buy
	} else {
		side = common.Sell
	}
	
	//check if we already have an open position in that dreiction
	positions, err := api.listPositions(q.GetSymbol())
	if err != nil {
		return fmt.Errorf("error getting positions: %v", err)
	}

	/**
	 * PS: We are only allowing one open trade at a time, we need to manually close the trade before opening another
	 * before opening another trade.
	 * This is strictly for testing purposes
	 */
	
	if len(positions) > 0 {
		return fmt.Errorf("position already open for %s", q.GetSymbol().Name)
	}

	return api.openPosition(q.GetSymbol().Name, 0.01, side)
}