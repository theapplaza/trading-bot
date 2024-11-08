package twelvedata

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiBaseUrl    string
	ApiKey        string
	DataStreamUrl string
	Instruments   string
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
		ApiBaseUrl:    os.Getenv("TWELVE_DATA_WS_URL"),
		ApiKey:        os.Getenv("TWELVE_DATA_API_KEY"),
		Instruments:   os.Getenv("INSTRUMENTS"),
		DataStreamUrl: os.Getenv("TWELVE_DATA_WS_URL"),
	}

}
