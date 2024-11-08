package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"trading-bot/core/handlers"

	"github.com/gorilla/websocket"
)

type EventType string

const (
	Subscribe EventType = "subscribe-status"
	Price     EventType = "price"
	MessageProcessing EventType = "message-processing"
)

type Respone struct {
	Event  EventType `json:"event"`
	Status string    `json:"status"`
}

type QuotePriceEvent struct {
	*Respone
	Symbol    string  `json:"symbol"`
	Currency  string  `json:"currency"`
	Exchange  string  `json:"exchange"`
	MicCode   string  `json:"mic_code"`
	Type      string  `json:"type"`
	Timestamp int     `json:"timestamp"`
	Price     float64 `json:"price"`
	DayVolume int     `json:"day_volume"`
}

var prices = make([]float64, 0)

func processPriceQuote(con *websocket.Conn, symbols []string) (err error){
	err = con.WriteMessage(websocket.TextMessage, []byte(`
		{
			"action": "subscribe",
			"params": {
				"symbols": "`+strings.Join(symbols, ",")+`",
			}
		}`,
	))
	if err != nil {
		return
	}

	// var data Respone
	// _, message, err := con.ReadMessage()
	// err = json.Unmarshal(message, &data)
	// if err != nil || data.Status == "error" {
	// 	log.Println("Error opening connection", err)
	// 	// log.Fatalf("Response %s  - %s", data.Event, data.Status)
	// }
		

	 // Read messages from the WebSocket
    for {
        _, message, err := con.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            return err
        }
        handleMessage(message)
    }
	
}

func handleMessage(msg []byte) {
	var data QuotePriceEvent
	err := json.Unmarshal(msg, &data)
	if err != nil {
		log.Println("Error unmarshaling WebSocket message:", err)
		return
	}

	if data.Event == Subscribe {
		log.Printf("Received values for %s", data.Event)
		return
	}

	if data.Event == Price {
		log.Printf("Received values for %s for %s with price %f", data.Event, data.Symbol, data.Price)
		prices = append(prices, data.Price)
 		handlers.Compute(prices)
		fmt.Println("Prices ", prices)
	}

}
