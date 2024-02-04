package logic

import (
	"errors"
	"math/rand"
)

func AddURL(m map[string]string, url string) string {
	for id := randomID(); true; id = randomID() {
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
		return "", errors.New("no URL with this ID")
	}
}

const (
	UppercaseA = 65
	LowercaseA = 97
	IdLen = 5
	LettersCount = 26
)

func randomID() string {
	var chars []byte
	for i := 0; i < IdLen; i++ {
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