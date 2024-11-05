package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/db"
	"github.com/iliamikado/UrlShortener/internal/grpc"
	"github.com/iliamikado/UrlShortener/internal/handlers"
	"github.com/iliamikado/UrlShortener/internal/logger"
	"github.com/iliamikado/UrlShortener/internal/storage"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	config.ParseConfig()

	if err := run(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("Server error: " + err.Error())
	}

	logger.Log.Info("Server Shutdown gracefully")
}

func run() error {
	if err := logger.Initialize(config.LoggerLevel); err != nil {
		return err
	}

	urlStorage := createStorageFromConfig()
	r := handlers.AppRouter(urlStorage)

	logger.Log.Info("Running server", zap.String("address", config.LaunchAddress))

	runDegugServer()
	go grpc.RunGRPCServer(":8088", urlStorage)

	var srv = createServer(config.LaunchAddress, r)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("HTTP server Shutdown. Errors: " + err.Error())
		}
	}()

	if config.EnableHTTPS {
		return srv.ListenAndServeTLS("cert.pem", "key.pem")
	} else {
		return srv.ListenAndServe()
	}
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

func runDegugServer() {
	go func() {
		http.ListenAndServe(config.DebugAddress, nil)
	}()
}

func createServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
