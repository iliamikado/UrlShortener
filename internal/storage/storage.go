package storage

import (
	"errors"
	"math/rand"

	"github.com/iliamikado/UrlShortener/internal/config"
)

type URLStorage struct {
	m map[string]string
	fileSaver *FileSaver
}

func NewURLStorage(pathToFile string) *URLStorage {
	var st URLStorage
	st.m = make(map[string]string)
	if pathToFile != "" {
		st.fileSaver = NewFileSaver(pathToFile)
		st.m = st.fileSaver.GetAllData()
	}
	return &st
}

func (st *URLStorage) AddURL(url string) string {
	for id := randomID(); true; id = randomID() {
		if _, exist := st.m[id]; !exist {
			st.m[id] = url
			if st.fileSaver != nil {
				st.fileSaver.AddURL(SavedURL{
					ID: id,
					ShortURL: config.ResultAddress + "/" + id,
					OriginURL: url,
				})
			}
			return id
		}
	}
	return ""
}

func (st *URLStorage) GetURL(id string) (string, error) {
	if url, ok := st.m[id]; ok {
		return url, nil
	} else {
		return "", errors.New("no URL with this ID")
	}
}

const (
	UppercaseA = 65
	LowercaseA = 97
	IDLen = 5
	LettersCount = 26
)

func randomID() string {
	var chars []byte
	for i := 0; i < IDLen; i++ {
		uppercase := rand.Intn(2)
		letter := rand.Intn(LettersCount)
		if uppercase == 0 {
			letter += LowercaseA
		} else {
			letter += UppercaseA
		}
		chars = append(chars, byte(letter))
	}
	return string(chars)
}