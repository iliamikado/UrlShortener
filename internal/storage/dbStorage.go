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

func (st *DBStorage) AddURL(longURL string, userID uint) (string, error) {
	var err *pgconn.PgError
	id, e := st.urlDB.AddURL(longURL, userID, randomID)
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

func (st *DBStorage) AddManyURLs(longURLs []string, userID uint) []string {
	ids, _ := st.urlDB.AddManyURLs(longURLs, userID, randomID)
	return ids
}

func (st *DBStorage) CreateNewUser() uint {
	return st.urlDB.CreateNewUser()
}

func (st *DBStorage) GetUserURLs(userID uint) [][2]string{
	return st.urlDB.GetUserURLs(userID)
}