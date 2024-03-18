package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/db"
	"github.com/iliamikado/UrlShortener/internal/handlers"
	"github.com/iliamikado/UrlShortener/internal/logger"
	"github.com/iliamikado/UrlShortener/internal/storage"
)

func main() {
	config.ParseConfig()

	if err := run(); err != nil {
		panic(err)
	}
}

var urlStorage *storage.URLStorage

func run() error {
	if err := logger.Initialize(config.LoggerLevel); err != nil {
        return err
    }

	db.Initialize(config.DatabaseDsn)

	urlStorage = storage.NewURLStorage(config.FileStoragePath)
	r := handlers.AppRouter(urlStorage)

	logger.Log.Info("Running server", zap.String("address", config.LaunchAddress))
	return http.ListenAndServe(config.LaunchAddress, r)
}
