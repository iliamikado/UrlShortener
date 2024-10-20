package storage

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/iliamikado/UrlShortener/internal/db"
)

// DBStorage - структура для работы с БД
type DBStorage struct {
	urlDB *db.URLShortenerDB
}

// Реализация URLStorage интерфейса
func NewDBStorage(urlDB *db.URLShortenerDB) *DBStorage {
	var st DBStorage
	urlDB.CreateURLTable()
	st.urlDB = urlDB
	return &st
}

// Реализация URLStorage интерфейса
func (st *DBStorage) AddURL(longURL string, userID string) (string, error) {
	var err *pgconn.PgError
	id, e := st.urlDB.AddURL(longURL, userID, randomID)
	errors.As(e, &err)
	if err != nil && err.Code == pgerrcode.UniqueViolation && err.ConstraintName == "urls_long_url_key" {
		id, _ = st.urlDB.GetIDByURL(longURL)
		return id, URLAlreadyExistsError
	}
	return id, err
}

// Реализация URLStorage интерфейса
func (st *DBStorage) GetURL(id string) (string, error) {
	url, isDeleted, err := st.urlDB.GetURL(id)
	if isDeleted {
		return url, URLIsDeleted
	}
	return url, err
}

// Реализация URLStorage интерфейса
func (st *DBStorage) AddManyURLs(longURLs []string, userID string) []string {
	ids, _ := st.urlDB.AddManyURLs(longURLs, userID, randomID)
	return ids
}

// Реализация URLStorage интерфейса
func (st *DBStorage) CreateNewUser() string {
	return uuid.NewString()
}

// Реализация URLStorage интерфейса
func (st *DBStorage) GetUserURLs(userID string) [][2]string {
	return st.urlDB.GetUserURLs(userID)
}

// Реализация URLStorage интерфейса
func (st *DBStorage) DeleteURLs(ids []string, userID string) {
	go st.urlDB.DeleteURLs(ids, userID)
}
