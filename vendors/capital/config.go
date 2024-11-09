package capital

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	ApiBaseUrl       string
	ApiKey           string
	ApiKeyPassword   string
	ApiKeyUser       string
	DataStreamUrl    string
	PriceUrl         string
	Instruments      string
	TradingAccountId string
}

var activeConfig *config

func init() {
	activeConfig = loadConfig()
}

func loadConfig() *config {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &config{
		ApiBaseUrl:       os.Getenv("CAPITAL_API_URL"),
		ApiKey:           os.Getenv("CAPITAL_API_KEY"),
		ApiKeyPassword:   os.Getenv("CAPITAL_API_KEY_PASSWORD"),
		ApiKeyUser:       os.Getenv("CAPITAL_IDENTIFIER"),
		TradingAccountId: os.Getenv("CAPITAL_PRIMARY_ACCOUNT_ID"),
		Instruments:      os.Getenv("CAPITAL_INSTRUMENTS"),
	}

}
