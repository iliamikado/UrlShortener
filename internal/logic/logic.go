package logic

import (
	"errors"
	"math/rand"
)

func AddURL(m map[string]string, url string) string {
	for id := randomId(); true; id = randomId() {
		if _, exist := m[id]; !exist {
			m[id] = url
			return id
		}
	}
	return ""
}

func GetURL(m map[string]string, id string) (string, error) {
	if url, ok := m[id]; ok {
		return url, nil
	} else {
		return "", errors.New("No URL with this ID")
	}
}

const (
	UPPERCASE_A = 65
	LOWERCASE_A = 97
	ID_LEN = 5
	LETTERS_COUNT = 26
)

func randomId() string {
	var chars []byte
	for i := 0; i < ID_LEN; i++ {
		uppercase := rand.Intn(2)
		letter := rand.Intn(LETTERS_COUNT)
		if uppercase == 0 {
			letter += LOWERCASE_A
		} else {
			letter += UPPERCASE_A
		}
		chars = append(chars, byte(letter))
	}
	return string(chars)
}