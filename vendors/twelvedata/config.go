package twelvedata

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	ApiBaseUrl    string
	ApiKey        string
	DataStreamUrl string
	Instruments   string
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
		ApiBaseUrl:    os.Getenv("TWELVE_DATA_WS_URL"),
		ApiKey:        os.Getenv("TWELVE_DATA_API_KEY"),
		Instruments:   os.Getenv("TWELVE_DATA_INSTRUMENTS"),
		DataStreamUrl: os.Getenv("TWELVE_DATA_WS_URL"),
	}

}
