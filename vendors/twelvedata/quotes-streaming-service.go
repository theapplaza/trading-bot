package twelvedata

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type QuoteStreamer struct{}

var errorChannel chan<- error 

func New(errorChan chan<- error) QuoteStreamer{
	errorChannel = errorChan
	return QuoteStreamer{}
}


// func streamQuotes()
func (QuoteStreamer) StreamQuotes() (err error) {

	defer func(){
		if err  != nil {
			errorChannel <- err
		}
	}()

	con, err := _createConnection()

	if err != nil {
		return
	}

	err = _subcribe(con)
	if err != nil {
		return
	}
	
	go _listen(con)

	return
}

func _createConnection() (con *websocket.Conn, err error) {
	url := fmt.Sprintf("wss://%v/quotes/price?apikey=%v", activeConfig.DataStreamUrl, activeConfig.ApiKey)
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
			"symbols": activeConfig.Instruments,
		},
	}

	// Send subscription request
	if err = con.WriteJSON(request); err != nil {
		return
	}

	return
}

func _listen(con *websocket.Conn) (err error) {

	defer func(){
		if err  != nil {
			errorChannel <- err
		}
		con.Close()
	}()
	

	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			err = fmt.Errorf("error reading message: %v", err)
			return err
		}

		var response map[string]interface{}
		if err := json.Unmarshal(message, &response); err != nil {
			errorChannel <- err
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
			return fmt.Errorf("unhandled message: %v", response)
		}


	}
}

func _handleSubscriptionResponse(response map[string]interface{}) {
	status := response["status"].(string)
	if status != "ok" {
		errorChannel <- fmt.Errorf("subscription failed: %v", response)
	}
}

func _handleQuoteUpdateResponse(payload map[string]interface{}) {
	epic := payload["symbol"].(string)
	price := payload["price"].(float64)
	timestamp := payload["timestamp"].(float64)

	log.Printf("TwelveData: Price update for %s - Price: %f, Timestamp: %d", epic, price, int64(timestamp))
}
