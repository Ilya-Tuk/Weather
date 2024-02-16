package config

import (
	"os"
)

type Config struct {
	ServCfg ServerConfig
	Db      DatabaseConfig
}

type ServerConfig struct {
	DebugMode     bool
	API_KEY       string
	WeatherApiUrl string
	ServerHost    string
}

type DatabaseConfig struct {
	MemoryFileName string
}

func setStringWithDefValue(variable *string, def string, envname string) {
	if os.Getenv(envname) == "" {
		*variable = def
	} else {
		*variable = os.Getenv(envname)
	}
}

func Read() Config {
	var debugMode bool
	var API_KEY string
	var weatherApiUrl string
	var serverHost string
	var memoryFileName string

	if os.Getenv("DEBUG") == "true" {
		debugMode = true
	}

	setStringWithDefValue(&serverHost, ":8080", "HOST")

	API_KEY = os.Getenv("API_KEY")

	setStringWithDefValue(&weatherApiUrl, "http://api.weatherapi.com/v1/", "API_URL")

	setStringWithDefValue(&memoryFileName, "db.txt", "DB_FILE")

	return Config{
		ServCfg: ServerConfig{
			DebugMode:     debugMode,
			ServerHost:    serverHost,
			API_KEY:       API_KEY,
			WeatherApiUrl: weatherApiUrl,
		},
		Db: DatabaseConfig{
			MemoryFileName: memoryFileName,
		},
	}
}
