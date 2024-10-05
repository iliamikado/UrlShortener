package storage

import "errors"

type URLStorage interface {
	AddURL(url string, userID string) (string, error)
	GetURL(id string) (string, error)
	AddManyURLs(longURLs []string, userID string) []string
	CreateNewUser() string
	GetUserURLs(userID string) [][2]string
	DeleteURLs(ids []string, userID string)
}

var (
	URLAlreadyExistsError error
	URLIsDeleted          error
)

func init() {
	URLAlreadyExistsError = errors.New("trying to add existing url")
	URLIsDeleted = errors.New("url has been deleted")
}
