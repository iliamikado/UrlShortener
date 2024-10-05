package storage

import "math/rand"

const (
	UppercaseA   = 65
	LowercaseA   = 97
	IDLen        = 5
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
