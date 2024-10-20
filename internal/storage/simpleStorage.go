package storage

import (
	"errors"

	"github.com/google/uuid"
)

// Структура для простого хранилища (в ОЗУ)
type SimpleStorage struct {
	m         map[string]string
	usersURLs map[string][]string
}

// Создание простого хранилища
func NewSimpleStorage() *SimpleStorage {
	var st SimpleStorage
	st.m = make(map[string]string)
	st.usersURLs = make(map[string][]string)
	return &st
}

// Реализация URLStorage интерфейса
func (st *SimpleStorage) AddURL(longURL string, userID string) (string, error) {
	var newID string
	for id := randomID(); newID == ""; id = randomID() {
		if _, exist := st.m[id]; !exist {
			newID = id
		}
	}
	st.m[newID] = longURL
	st.usersURLs[userID] = append(st.usersURLs[userID], newID)
	return newID, nil
}

// Реализация URLStorage интерфейса
func (st *SimpleStorage) GetURL(id string) (string, error) {
	if url, ok := st.m[id]; ok {
		return url, nil
	} else {
		return "", errors.New("no URL with this ID")
	}
}

// Реализация URLStorage интерфейса
func (st *SimpleStorage) AddManyURLs(longURLs []string, userID string) []string {
	var ids []string
	for _, url := range longURLs {
		id, _ := st.AddURL(url, userID)
		ids = append(ids, id)
	}
	return ids
}

// Реализация URLStorage интерфейса
func (st *SimpleStorage) CreateNewUser() string {
	return uuid.NewString()
}

// Реализация URLStorage интерфейса
func (st *SimpleStorage) GetUserURLs(userID string) [][2]string {
	var ans = make([][2]string, len(st.usersURLs[userID]))
	for i, urlID := range st.usersURLs[userID] {
		ans[i] = [2]string{urlID, st.m[urlID]}
	}
	return ans
}

// Реализация URLStorage интерфейса
func (st *SimpleStorage) DeleteURLs(ids []string, userID string) {
	// TODO
}
