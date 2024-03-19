package db

import (
	"context"
	"database/sql"

	"github.com/iliamikado/UrlShortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type URLShortenerDB struct {
	db *sql.DB
}
var URLDB URLShortenerDB

func Initialize(host string) {
	logger.Log.Info("Opening DB with host=" + host)
	db, err := sql.Open("pgx", host)
	if err != nil {
		logger.Log.Error("Failed to open DB")
		panic(err)
	}
	URLDB = URLShortenerDB{db}
}

func (urlDB *URLShortenerDB) CreateURLTable() {
	urlDB.db.Exec(`create table if not exists urls (
		id text PRIMARY KEY NOT NULL,
		long_url text UNIQUE NOT NULL
	);`)
}

func (urlDB *URLShortenerDB) AddURL(id, longURL string) error {
	_, err := urlDB.db.Exec(`insert into urls values ($1, $2)`, id, longURL)
	return err
}

func (urlDB *URLShortenerDB) AddManyURLs(ids, longURLs []string) error {
	tx, _ := urlDB.db.Begin()
	stmt, _ := tx.Prepare(`insert into urls values ($1, $2)`)
	defer stmt.Close()
	for i := 0; i < len(ids); i++ {
		stmt.Exec(ids[i], longURLs[i])
	}
	return tx.Commit()
}

func (urlDB *URLShortenerDB) GetURL(id string) (string, error) {
	row := urlDB.db.QueryRow(`select long_url from urls where id = $1`, id)
	var longURL string
	err := row.Scan(&longURL)
	return longURL, err
}

func (urlDB *URLShortenerDB) GetIDByURL(longURL string) (string, error) {
	row := urlDB.db.QueryRow(`select id from urls where long_url = $1`, longURL)
	var id string
	err := row.Scan(&id)
	return id, err
}

func (urlDB *URLShortenerDB) Close() {
	urlDB.db.Close()
}

func (urlDB *URLShortenerDB) Ping() error {
	return urlDB.db.PingContext(context.TODO())
}