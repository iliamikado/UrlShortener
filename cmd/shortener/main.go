package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"

	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/db"
	"github.com/iliamikado/UrlShortener/internal/handlers"
	"github.com/iliamikado/UrlShortener/internal/logger"
	"github.com/iliamikado/UrlShortener/internal/storage"
)

var (
	buildVersion string = "N/A"
	buildDate string = "N/A"
	buildCommit string = "N/A"
)

func main() {

	fmt.Println("Build version: " + buildVersion)
	fmt.Println("Build date: " + buildDate)
	fmt.Println("Build commit: " + buildCommit)

	config.ParseConfig()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(config.LoggerLevel); err != nil {
		return err
	}

	urlStorage := createStorageFromConfig()
	r := handlers.AppRouter(urlStorage)

	logger.Log.Info("Running server", zap.String("address", config.LaunchAddress))
	go func() {
		http.ListenAndServe(config.DebugAddress, nil)
	}()
	return http.ListenAndServe(config.LaunchAddress, r)
}

func createStorageFromConfig() storage.URLStorage {
	if config.DatabaseDsn != "" {
		db.Initialize(config.DatabaseDsn)
		return storage.NewDBStorage(&db.URLDB)
	} else if config.FileStoragePath != "" {
		return storage.NewDiskStorage(config.FileStoragePath)
	} else {
		return storage.NewSimpleStorage()
	}
}
