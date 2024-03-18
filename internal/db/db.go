package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iliamikado/UrlShortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UrlShortenerDb struct {
	db *sql.DB
}
var UrlDb UrlShortenerDb

func Initialize(host string) {
	logger.Log.Info("Opening DB with host=" + host)
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, `admin`, `password`, `url-shortener`)
	db, err := sql.Open("pgx", ps)
	if err != nil {
		logger.Log.Error("Failed to open DB")
		panic(err)
	}
	UrlDb = UrlShortenerDb{db}
}

func (urlDb *UrlShortenerDb) Close() {
	urlDb.db.Close()
}

func (urlDb *UrlShortenerDb) Ping() error {
	return urlDb.db.PingContext(context.TODO())
}