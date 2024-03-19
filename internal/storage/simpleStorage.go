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
	var newId string
	for id := randomID(); newId == ""; id = randomID() {
		if _, exist := st.m[id]; !exist {
			newId = id
		}
	}
	st.m[newId] = longURL
	return newId
}

func (st *SimpleStorage) GetURL(id string) (string, error) {
	if url, ok := st.m[id]; ok {
		return url, nil
	} else {
		return "", errors.New("no URL with this ID")
	}
}