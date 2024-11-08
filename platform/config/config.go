package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiBaseUrl       string
	ApiKey           string
	ApiKeyPassword   string
	ApiKeyUser       string
	WebSocketURl     string
	PriceUrl         string
	Instruments      string
	TradingAccountId string
}

func LoadConfig() Config {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//@TODO find a way to know the vendor info to use from runtime env or params
	client := os.Getenv("ACTIVE_CLIENT")
	var config *Config
	switch client {
	case "CAPITAL":
		config = loadCaptalConfig()
	default:
		config = loadTwelveDataConfig()
	}

	config.Instruments = os.Getenv("INSTRUMENTS")

	return *config

}

func loadCaptalConfig() *Config {
	return &Config{
		ApiBaseUrl:       os.Getenv("CAPITAL_API_URL"),
		ApiKey:           os.Getenv("CAPITAL_API_KEY"),
		ApiKeyPassword:   os.Getenv("CAPITAL_API_KEY_PASSWORD"),
		ApiKeyUser:       os.Getenv("CAPITAL_IDENTIFIER"),
		TradingAccountId: os.Getenv("CAPITAL_PRIMARY_ACCOUNT_ID"),
	}
}

func loadTwelveDataConfig() *Config {
	return &Config{
		ApiBaseUrl:   os.Getenv("TWELVE_DATA_API_URL"),
		WebSocketURl: os.Getenv("TWELVE_DATA_WS_URL"),
		ApiKey:       os.Getenv("TWELVE_DATA_API_KEY"),
		PriceUrl:     os.Getenv("TWELVE_DATA_PRICE_URL"),
	}
}
