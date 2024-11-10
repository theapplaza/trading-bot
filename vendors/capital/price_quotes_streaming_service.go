package capital

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gorilla/websocket"
)

type QuoteStreamer struct {
}

func New() QuoteStreamer {
	return QuoteStreamer{}
}

func(s QuoteStreamer) GetName() string {
	return "Capital"
}

func (s QuoteStreamer) StreamQuotes() (err error) {

	//ensure that authentication is done
	if activeSession == nil {
		if err = authenticate(); err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s%s", activeSession.StreamingHost, "connect")
	con, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer con.Close()

	//subcribe to price update
	request := map[string]interface{}{
		"destination":   "marketData.subscribe",
		"correlationId": rand.Int(),
		"cst":           activeSession.cst,
		"securityToken": activeSession.securitytoken,
		"payload": map[string]interface{}{
			"epics": strings.Split(activeConfig.Instruments, ","),
		},
	}
	err = con.WriteJSON(request)

	return s.listen(con)
}

func (s QuoteStreamer) listen(con *websocket.Conn) (err error) {

	defer con.Close()

	for {
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
		case "marketData.subscribe":
			s.handleSubscriptionResponse(response)
		case "quote":
			s.handleQuoteUpdateResponse(response)
		default:
			return fmt.Errorf("unhandled message: %v", response)
		}
	}
}

func (s QuoteStreamer) handleSubscriptionResponse(response map[string]interface{}) error {
	status := response["status"].(string)
	if status != "OK" {
		return fmt.Errorf("subscription failed: %v", response)
	}
	return nil
}

func (s QuoteStreamer) handleQuoteUpdateResponse(response map[string]interface{}) {
	payload := response["payload"].(map[string]interface{})

	epic := payload["epic"].(string)
	bid := payload["bid"].(float64)
	ofr := payload["ofr"].(float64)
	timestamp := payload["timestamp"].(float64)

	log.Printf("Capital: Price update for %s - Bid: %f, Offer: %.2f, Timestamp: %d", epic, bid, ofr, int64(timestamp))
}
