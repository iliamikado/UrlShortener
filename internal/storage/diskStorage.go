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
	st.smSt.m = st.fileSaver.GetAllData()
	return &st
}

func (st *DiskStorage) AddURL(longURL string) (string, error) {
	id, err := st.smSt.AddURL(longURL)
	st.fileSaver.AddURL(SavedURL{
		ID:        id,
		ShortURL:  config.ResultAddress + "/" + id,
		OriginURL: longURL,
	})
	return id, err
}

func (st *DiskStorage) GetURL(id string) (string, error) {
	return st.smSt.GetURL(id)
}

func (st *DiskStorage) AddManyURLs(longURLs []string) []string {
	var ids []string
	for _, url := range longURLs {
		id, _ := st.AddURL(url)
		ids = append(ids, id)
	}
	return ids
}