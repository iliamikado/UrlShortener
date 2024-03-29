package storage

import (
	"errors"

	"github.com/iliamikado/UrlShortener/internal/db"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBStorage struct {
	urlDB *db.URLShortenerDB
}

func NewDBStorage(urlDB *db.URLShortenerDB) *DBStorage {
	var st DBStorage
	urlDB.CreateURLTable()
	st.urlDB = urlDB
	return &st
}

func (st *DBStorage) AddURL(longURL string) (string, error) {
	var err *pgconn.PgError
	id, e := st.urlDB.AddURL(longURL, randomID)
	errors.As(e, &err)
	if err != nil && err.Code == pgerrcode.UniqueViolation && err.ConstraintName == "urls_long_url_key" {
		id, _ = st.urlDB.GetIDByURL(longURL)
		return id, URLAlreadyExistsError
	}
	return id, err
}

func (st *DBStorage) GetURL(id string) (string, error) {
	return st.urlDB.GetURL(id)
}

func (st *DBStorage) AddManyURLs(longURLs []string) []string {
	ids, _ := st.urlDB.AddManyURLs(longURLs, randomID)
	return ids
}