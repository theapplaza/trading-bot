package core

import (
	"fmt"
	"log"
	"trading-bot/platform/config"

	"github.com/gorilla/websocket"
)

//Subscribes to fetching prices which we can use for other information
func Start(cfg config.Config) {

	con, err := setupConnection(cfg)
	if err != nil {
        log.Fatal("WebSocket connection failed:", err)
	}
	defer con.Close();

			
	err = processPriceQuote(con, []string{"EUR/USD"})
	if err != nil {
        log.Fatal("WebSocket subscription failed:", err)
	}
}

func setupConnection(cfg config.Config) (con *websocket.Conn, err error) {
	url := fmt.Sprintf("wss://%v/%v?apikey=%v", cfg.WebSocketURl, cfg.PriceUrl, cfg.ApiKey)
	con, _, err = websocket.DefaultDialer.Dial(url, nil)
	return
}

