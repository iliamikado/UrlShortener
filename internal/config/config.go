package config

import (
	"flag"
	"os"
)

// Переменные окружения и флаги
var (
	// LaunchAddress - адрес запуска сервера
	LaunchAddress string
	// ResultAddress - адрес результата для коротких ссылок
	ResultAddress string
	// LoggerLevel - уровень логгера
	LoggerLevel string
	// FileStoragePath - путь до файла с сохранением
	FileStoragePath string
	// DatabaseDsn - строка для подключения к бд
	DatabaseDsn string
	// DebugAddress - адрес дебага
	DebugAddress string
	// EnableHttps - включать ли https
	EnableHttps bool
)

// ParseConfig - чтение конфига из флагов и переменных окружения
func ParseConfig() {
	flag.StringVar(&LaunchAddress, "a", "localhost:8080", "Set launch address for server")
	flag.StringVar(&ResultAddress, "b", "http://localhost:8080", "Set basic result address for short URL")
	flag.StringVar(&LoggerLevel, "l", "info", "Set Logger level")
	flag.StringVar(&FileStoragePath, "f", "/tmp/short-url-db.json", "Set file storage path for urls")
	flag.StringVar(&DatabaseDsn, "d", "", "Set DB adress")
	flag.StringVar(&DebugAddress, "g", "localhost:8081", "Set debug address")
	flag.BoolVar(&EnableHttps, "s", false, "Enable https")
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
	if enableHttps := os.Getenv("ENABLE_HTTPS"); enableHttps == "true" {
		EnableHttps = true
	}
}
