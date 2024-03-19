package storage

import "errors"

type SimpleStorage struct {
	m map[string]string
}

func NewSimpleStorage() *SimpleStorage {
	var st SimpleStorage
	st.m = make(map[string]string)
	return &st
}

func (st *SimpleStorage) AddURL(longURL string) string {
	var newID string
	for id := randomID(); newID == ""; id = randomID() {
		if _, exist := st.m[id]; !exist {
			newID = id
		}
	}
	st.m[newID] = longURL
	return newID
}

func (st *SimpleStorage) GetURL(id string) (string, error) {
	if url, ok := st.m[id]; ok {
		return url, nil
	} else {
		return "", errors.New("no URL with this ID")
	}
}

func (st *SimpleStorage) AddManyURLs(longURLs []string) []string {
	var ids []string
	for _, url := range longURLs {
		ids = append(ids, st.AddURL(url))
	}
	return ids
}