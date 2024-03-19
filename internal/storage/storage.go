package storage

type URLStorage interface {
	AddURL(url string) string
	GetURL(id string) (string, error)
}
