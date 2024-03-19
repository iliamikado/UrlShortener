package storage

import "github.com/iliamikado/UrlShortener/internal/db"

type DBStorage struct {
	urlDB *db.URLShortenerDB
}

func NewDBStorage(urlDB *db.URLShortenerDB) *DBStorage {
	var st DBStorage
	urlDB.CreateURLTable()
	st.urlDB = urlDB
	return &st
}

func (st *DBStorage) AddURL(longURL string) string {
	var id string
	for id = randomID(); true; id = randomID() {
		err := st.urlDB.AddURL(id, longURL)
		if err == nil {
			break;
		}
	}	
	return id
}

func (st *DBStorage) GetURL(id string) (string, error) {
	return st.urlDB.GetURL(id)
}