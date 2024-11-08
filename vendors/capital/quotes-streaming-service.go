package capital

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

type destinatation string
const (
	quoteRequest = "marketData.subscribe"
	quote = "quote"
)

type subscriptionResponse struct {
	Event  string `json:"event"`
	Status string `json:"status"`
}

type QuoteStreamer struct{}

func New() *QuoteStreamer{
	return &QuoteStreamer{}
}

// func streamQuotes()
func (*QuoteStreamer) StreamQuotes()(err error) {

	//ensure that authentication is done
	if activeSession == nil {
		authenticate()
	}

	con, err := _createConnection()

	if err != nil {
		return
	}

	defer con.Close()

	err = _subcribe(con)
	if err != nil {
		return
	}

	_listen(con)

	return
}

func _createConnection() (con *websocket.Conn, err error) {
	url := fmt.Sprintf("%s%s", activeSession.StreamingHost, "connect")
	con, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		err = fmt.Errorf("error connecting to data streaming endpoint %s", err)
	}
	return
}

func _subcribe(con *websocket.Conn) (err error) {
	//subcribe to price update
	request := map[string]interface{}{
		"destination":   "marketData.subscribe",
		"correlationId": rand.Int(),
		"cst":           activeSession.cst,
		"securityToken": activeSession.securitytoken,
		"payload": map[string]interface{}{
			"epics": []string{"OIL_CRUDE"},
		},
	}

	// Send subscription request
	err = con.WriteJSON(request); 
	return
}

func _listen(con *websocket.Conn) (err error) {
	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
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
			_handleSubscriptionResponse(response)
		case "quote":
			_handleQuoteUpdateResponse(response)
		default:
			log.Println("Unhandled message:", response)
		}
	}
}

func _handleSubscriptionResponse(response map[string]interface{}) {
	status := response["status"].(string)
	if status != "OK" {
		log.Println("Subscription failed:", response)
	}
}

func _handleQuoteUpdateResponse(response map[string]interface{}) {
    payload := response["payload"].(map[string]interface{})

    epic := payload["epic"].(string)
    bid := payload["bid"].(float64)
    ofr := payload["ofr"].(float64)
    timestamp := payload["timestamp"].(float64)

    log.Printf("Capital: Price update for %s - Bid: %f, Offer: %.2f, Timestamp: %d", epic, bid, ofr, int64(timestamp))
}