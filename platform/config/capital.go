package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type CapitalConfig struct {
	WebSocketURl string
	APIKey       string
	PriceUrl     string
	INSTRUMENTS  string
}

func LoadCapitalConfig() CapitalConfig {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//@TODO find a way to know the vendor info to use from runtime env or params
	return CapitalConfig{
		os.Getenv("TWELVE_DATA_WS_URL"),
		os.Getenv("TWELVE_DATA_API_KEY"),
		os.Getenv("TWELVE_DATA_PRICE_URL"),
		os.Getenv("INSTRUMENTS"),
	}
}


