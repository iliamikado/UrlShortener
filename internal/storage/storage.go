package storage

import "errors"

type URLStorage interface {
	AddURL(url string) (string, error)
	GetURL(id string) (string, error)
	AddManyURLs(longURLs []string) []string
}

var URLAlreadyExistsError error

func init() {
	URLAlreadyExistsError = errors.New("trying to add existing url")
}
