package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iliamikado/UrlShortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type URLShortenerDB struct {
	db *sql.DB
}
var URLDB URLShortenerDB

func Initialize(host string) {
	logger.Log.Info("Opening DB with host=" + host)
	ps := fmt.Sprintf("host=%s", host)
	db, err := sql.Open("pgx", ps)
	if err != nil {
		logger.Log.Error("Failed to open DB")
		panic(err)
	}
	URLDB = URLShortenerDB{db}
}

func (urlDB *URLShortenerDB) Close() {
	urlDB.db.Close()
}

func (urlDB *URLShortenerDB) Ping() error {
	return urlDB.db.PingContext(context.TODO())
}