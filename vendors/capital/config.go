package capital

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
	DataStreamUrl    string
	PriceUrl         string
	Instruments      string
	TradingAccountId string
}

var config *Config

func init() {
	config = loadConfig()
}

func loadConfig() *Config {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		ApiBaseUrl:       os.Getenv("CAPITAL_API_URL"),
		ApiKey:           os.Getenv("CAPITAL_API_KEY"),
		ApiKeyPassword:   os.Getenv("CAPITAL_API_KEY_PASSWORD"),
		ApiKeyUser:       os.Getenv("CAPITAL_IDENTIFIER"),
		TradingAccountId: os.Getenv("CAPITAL_PRIMARY_ACCOUNT_ID"),
		Instruments:      os.Getenv("INSTRUMENTS"),
	}

}
