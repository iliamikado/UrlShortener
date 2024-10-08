package config

import (
	"flag"
	"os"
)

var (
	LaunchAddress   string
	ResultAddress   string
	LoggerLevel     string
	FileStoragePath string
	DatabaseDsn     string
	DebugAddress    string
)

func ParseConfig() {
	flag.StringVar(&LaunchAddress, "a", "localhost:8080", "Set launch address for server")
	flag.StringVar(&ResultAddress, "b", "http://localhost:8080", "Set basic result address for short URL")
	flag.StringVar(&LoggerLevel, "l", "info", "Set Logger level")
	flag.StringVar(&FileStoragePath, "f", "/tmp/short-url-db.json", "Set file storage path for urls")
	flag.StringVar(&DatabaseDsn, "d", "", "Set DB adress")
	flag.StringVar(&DebugAddress, "g", "localhost:8081", "Set debug address")
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
	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		FileStoragePath = fileStoragePath
	}
	if databaseDsn := os.Getenv("DATABASE_DSN"); databaseDsn != "" {
		DatabaseDsn = databaseDsn
	}
}
