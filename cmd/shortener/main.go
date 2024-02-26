package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iliamikado/UrlShortener/internal/storage"
	"github.com/iliamikado/UrlShortener/internal/config"
)

func main() {
	config.ParseConfig()

	if err := run(); err != nil {
		panic(err)
	}
}

var urlStorage *storage.URLStorage
func run() error {
	urlStorage = storage.NewURLStorage()
	r := AppRouter()
	return http.ListenAndServe(config.LaunchAddress, r)
}

func AppRouter() *chi.Mux{
	r := chi.NewRouter()
	r.Get("/{id}", getURL)
	r.Post("/", postURL)
	return r
}

func postURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	longURL := string(body)
	id := urlStorage.AddURL(longURL)
	shortURL := config.ResultAddress + "/" + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func getURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	longURL, err := urlStorage.GetURL(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}


