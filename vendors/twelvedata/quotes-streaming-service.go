package twelvedata

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type QuoteStreamer struct{}

func New() *QuoteStreamer{
	return &QuoteStreamer{}
}

// func streamQuotes()
func (*QuoteStreamer) StreamQuotes() (err error) {
	con, err := _createConnection()

	if err != nil {
		return
	}

	// defer con.Close()

	err = _subcribe(con)
	if err != nil {
		return
	}

	go _listen(con)

	// err = _listen(con)
	// if err != nil {
	// 	log.Printf("we stopped listening in twelve data because of %s", err)
	// }

	return
}

func _createConnection() (con *websocket.Conn, err error) {
	url := fmt.Sprintf("wss://%v/quotes/price?apikey=%v", config.DataStreamUrl, config.ApiKey)
	con, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		err = fmt.Errorf("error connecting to data streaming endpoint %s", err)
	}
	return
}

func _subcribe(con *websocket.Conn) (err error) {
	//subcribe to price update
	request := map[string]interface{}{
		"action": "subscribe",
		"params": map[string]interface{}{
			"symbols": config.Instruments,
		},
	}

	// Send subscription request
	if err = con.WriteJSON(request); err != nil {
		return
	}

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
		switch response["event"] {
		case "connection":
			return fmt.Errorf("issue with connection request: %s", response["messages"])
		case "subscribe-status":
			_handleSubscriptionResponse(response)
		case "price":
			_handleQuoteUpdateResponse(response)
		default:
			log.Println("Unhandled message:", response)
		}
	}
}

func _handleSubscriptionResponse(response map[string]interface{}) {
	status := response["status"].(string)
	if status != "ok" {
		log.Println("Subscription failed:", response)
	}
}

func _handleQuoteUpdateResponse(payload map[string]interface{}) {
	epic := payload["symbol"].(string)
	price := payload["price"].(float64)
	timestamp := payload["timestamp"].(float64)

	log.Printf("TwelveData: Price update for %s - Price: %f, Timestamp: %d", epic, price, int64(timestamp))
}
