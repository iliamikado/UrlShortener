package config

import (
	"flag"
	"os"
)

var LaunchAddress string
var ResultAddress string

func ParseConfig() {
	flag.StringVar(&LaunchAddress, "a", "localhost:8080", "Set launch address for server")
	flag.StringVar(&ResultAddress, "b", "http://localhost:8080", "Set basic result address for short URL")
	flag.Parse()

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		LaunchAddress = serverAddress
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		ResultAddress = baseURL
	}
}