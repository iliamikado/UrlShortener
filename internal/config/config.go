package config

import (
	"flag"
	"os"
)

var (
	LaunchAddress string
	ResultAddress string
	LoggerLevel string
)

func ParseConfig() {
	flag.StringVar(&LaunchAddress, "a", "localhost:8080", "Set launch address for server")
	flag.StringVar(&ResultAddress, "b", "http://localhost:8080", "Set basic result address for short URL")
	flag.StringVar(&LoggerLevel, "l", "info", "Set Logger level")
	flag.Parse()

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		LaunchAddress = serverAddress
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		ResultAddress = baseURL
	}
	if loggerLevel := os.Getenv("LOGGER_LEVEL"); loggerLevel != "" {
		LoggerLevel = loggerLevel
	}
}