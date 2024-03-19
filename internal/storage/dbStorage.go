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
	var id string
	var err *pgconn.PgError
	for id = randomID(); true; id = randomID() {
		e := st.urlDB.AddURL(id, longURL)
		if errors.As(e, &err) {
			if err.Code == pgerrcode.UniqueViolation && err.ConstraintName == "urls_id_key" {
				continue
			} else {
				break
			}
		}
	}
	if err.Code == pgerrcode.UniqueViolation && err.ConstraintName == "urls_long_url_key" {
		id, _ = st.urlDB.GetIDByURL(longURL)
		return id, URLAlreadyExistsError
	}
	return id, err
}

func (st *DBStorage) GetURL(id string) (string, error) {
	return st.urlDB.GetURL(id)
}

func (st *DBStorage) AddManyURLs(longURLs []string) []string {
	var ids []string
	for {
		ids = make([]string, len(longURLs))
		for i := 0; i < len(longURLs); i++ {
			ids[i] = randomID();
		}
		err := st.urlDB.AddManyURLs(ids, longURLs)
		if err == nil {
			break
		}
	}
	return ids
}