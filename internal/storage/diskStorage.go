package storage

import "github.com/iliamikado/UrlShortener/internal/config"

type DiskStorage struct {
	smSt      *SimpleStorage
	fileSaver *FileSaver
}

func NewDiskStorage(pathToFile string) *DiskStorage {
	var st DiskStorage
	st.smSt = NewSimpleStorage()
	st.fileSaver = NewFileSaver(pathToFile)
	st.smSt.m, st.smSt.usersURLs = st.fileSaver.GetAllData()
	return &st
}

func (st *DiskStorage) AddURL(longURL string, userID string) (string, error) {
	id, err := st.smSt.AddURL(longURL, userID)
	st.fileSaver.AddURL(SavedURL{
		ID:        id,
		ShortURL:  config.ResultAddress + "/" + id,
		OriginURL: longURL,
		UserID:    userID,
	})
	return id, err
}

func (st *DiskStorage) GetURL(id string) (string, error) {
	return st.smSt.GetURL(id)
}

func (st *DiskStorage) AddManyURLs(longURLs []string, userID string) []string {
	var ids []string
	for _, url := range longURLs {
		id, _ := st.AddURL(url, userID)
		ids = append(ids, id)
	}
	return ids
}

func (st *DiskStorage) CreateNewUser() string {
	return st.smSt.CreateNewUser()
}

func (st *DiskStorage) GetUserURLs(userID string) [][2]string {
	return st.smSt.GetUserURLs(userID)
}

func (st *DiskStorage) DeleteURLs(ids []string, userID string) {
	st.smSt.DeleteURLs(ids, userID)
}
