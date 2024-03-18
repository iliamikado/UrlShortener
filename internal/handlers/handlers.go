package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/db"
	"github.com/iliamikado/UrlShortener/internal/logger"
	"github.com/iliamikado/UrlShortener/internal/storage"
)

var urlStorage *storage.URLStorage

func AppRouter(st *storage.URLStorage) *chi.Mux {
	urlStorage = st

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(gzipMiddleware)
	r.Get("/{id}", getURL)
	r.Post("/", postURL)
	r.Post("/api/shorten", postJSON)
	r.Get("/ping", pingDB)
	return r
}

func createShortURL(longURL string) string {
	id := urlStorage.AddURL(longURL)
	shortURL := config.ResultAddress + "/" + id
	return shortURL
}

func postURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	longURL := string(body)
	shortURL := createShortURL(longURL)
	w.Header().Set("Content-Type", "text/plain")
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

type (
	RequestJSON struct {
		URL		string	`json:"url"`
	}
	ResponseJSON struct {
		Result	string	`json:"result"`
	}
)

func postJSON(w http.ResponseWriter, r *http.Request) {
	var (
		req RequestJSON
		buf bytes.Buffer
	)
	_, err := buf.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	longURL := req.URL
	shortURL := createShortURL(longURL)
	resp := ResponseJSON{shortURL}
	
	var body []byte
	if body, err = json.Marshal(resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func pingDB(w http.ResponseWriter, r *http.Request) {
	err := db.UrlDb.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
