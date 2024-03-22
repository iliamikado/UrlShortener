package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		uuid text PRIMARY KEY NOT NULL,
		short_url text NOT NULL,
		long_url text UNIQUE NOT NULL
	);`)
}

func (urlDB *URLShortenerDB) AddURL(longURL string, getId func() string) (string, error) {
	var err error
	var id string
	for {
		id = getId()
		shortURL := config.ResultAddress + "/" + id;
		_, err = urlDB.db.Exec(`insert into urls values ($1, $2, $3)`, id, shortURL, longURL)
		if err != nil && isErrorWithID(err) {
			continue
		}
		break
	}
	return id, err
}

func (urlDB *URLShortenerDB) AddManyURLs(longURLs []string, getId func() string) ([]string, error) {
	tx, _ := urlDB.db.Begin()
	stmt, _ := tx.Prepare(`insert into urls values ($1, $2, $3)`)
	defer stmt.Close()
	ids := make([]string, len(longURLs))
	for i := 0; i < len(longURLs); i++ {
		var err error;
		for {
			ids[i] = getId()
			shortURL := config.ResultAddress + "/" + ids[i];
			_, err = stmt.Exec(ids[i], shortURL, longURLs[i])
			if err != nil && isErrorWithID(err) {
				continue
			}
			break
		}
		if (err != nil) {
			return nil, tx.Rollback()
		}
	}
	return ids, tx.Commit()
}

func (urlDB *URLShortenerDB) GetURL(id string) (string, error) {
	row := urlDB.db.QueryRow(`select long_url from urls where uuid = $1`, id)
	var longURL string
	err := row.Scan(&longURL)
	return longURL, err
}

func (urlDB *URLShortenerDB) GetIDByURL(longURL string) (string, error) {
	row := urlDB.db.QueryRow(`select uuid from urls where long_url = $1`, longURL)
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

func isErrorWithID(e error) bool {
	var err *pgconn.PgError
	if errors.As(e, &err) {
		if err.Code == pgerrcode.UniqueViolation && err.ConstraintName == "urls_uuid_key" {
			return true
		}
	}
	return false
}