package config

import (
	"encoding/json"
	"flag"
	"io"
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
	// EnableHTTPS - включать ли https
	EnableHTTPS bool
	// ConfigFile - название файла с конфигом
	ConfigFile string
)

// ConfigJSON - формат конфиг файла
type ConfigJSON struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDsn     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// ParseConfig - чтение конфига из флагов и переменных окружения
func ParseConfig() {
	flag.StringVar(&LaunchAddress, "a", "localhost:8080", "Set launch address for server")
	flag.StringVar(&ResultAddress, "b", "http://localhost:8080", "Set basic result address for short URL")
	flag.StringVar(&LoggerLevel, "l", "info", "Set Logger level")
	flag.StringVar(&FileStoragePath, "f", "/tmp/short-url-db.json", "Set file storage path for urls")
	flag.StringVar(&DatabaseDsn, "d", "", "Set DB adress")
	flag.StringVar(&DebugAddress, "g", "localhost:8081", "Set debug address")
	flag.BoolVar(&EnableHTTPS, "s", false, "Enable https")
	flag.StringVar(&ConfigFile, "c", "", "Set config file")
	flag.Parse()

	if configFile := os.Getenv("CONFIG"); configFile != "" {
		ConfigFile = configFile
	}

	if ConfigFile != "" {
		file, _ := os.Open(ConfigFile)
		b, _ := io.ReadAll(file)
		var configData ConfigJSON
		json.Unmarshal(b, &configData)
		if flag.Lookup("a") == nil && configData.ServerAddress != "" {
			LaunchAddress = configData.ServerAddress
		}
		if flag.Lookup("b") == nil && configData.BaseURL != "" {
			ResultAddress = configData.BaseURL
		}
		if flag.Lookup("f") == nil && configData.FileStoragePath != "" {
			FileStoragePath = configData.FileStoragePath
		}
		if flag.Lookup("d") == nil && configData.DatabaseDsn != "" {
			DatabaseDsn = configData.DatabaseDsn
		}
		if flag.Lookup("s") == nil {
			EnableHTTPS = configData.EnableHTTPS
		}
	}

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
		EnableHTTPS = true
	}
}
