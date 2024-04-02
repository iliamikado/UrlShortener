package storage

import "errors"

type URLStorage interface {
	AddURL(url string, userID uint) (string, error)
	GetURL(id string) (string, error)
	AddManyURLs(longURLs []string, userID uint) []string
	CreateNewUser() uint
	GetUserURLs(userID uint) [][2]string
	DeleteURLs(ids []string, userID uint)
}

var (
	URLAlreadyExistsError error
	URLIsDeleted error
)

func init() {
	URLAlreadyExistsError = errors.New("trying to add existing url")
	URLIsDeleted = errors.New("url has been deleted")
}
