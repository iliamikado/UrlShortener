package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/db"
	"github.com/iliamikado/UrlShortener/internal/logger"
	"github.com/iliamikado/UrlShortener/internal/storage"
)

var urlStorage storage.URLStorage

func AppRouter(st storage.URLStorage) *chi.Mux {
	urlStorage = st

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(gzipMiddleware)
	r.Get("/{id}", getURL)
	r.Post("/", authMiddleware(postURL))
	r.Post("/api/shorten", authMiddleware(postJSON))
	r.Get("/ping", pingDB)
	r.Post("/api/shorten/batch", authMiddleware(postManyURL))
	r.Get("/api/user/urls", getUserURLs)
	r.Delete("/api/user/urls", authMiddleware(deleteURLs))

	return r
}

func createShortURL(longURL string, userID string) (string, error) {
	logger.Log.Info("Get longUrl = " + longURL)
	id, err := urlStorage.AddURL(longURL, userID)
	shortURL := config.ResultAddress + "/" + id
	return shortURL, err
}

func postURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	longURL := string(body)
	userID := r.Context().Value(userIDKey{}).(string)
	shortURL, err := createShortURL(longURL, userID)
	w.Header().Set("Content-Type", "text/plain")
	if err != nil && errors.Is(err, storage.URLAlreadyExistsError) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(shortURL))
}

func getURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	longURL, err := urlStorage.GetURL(id)
	if errors.Is(err, storage.URLIsDeleted) {
		w.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type (
	RequestJSON struct {
		URL string `json:"url"`
	}
	ResponseJSON struct {
		Result string `json:"result"`
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
	userID := r.Context().Value(userIDKey{}).(string)
	shortURL, createErr := createShortURL(longURL, userID)
	resp := ResponseJSON{shortURL}

	var body []byte
	if body, err = json.Marshal(resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if createErr != nil && errors.Is(createErr, storage.URLAlreadyExistsError) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(body)
}

func pingDB(w http.ResponseWriter, r *http.Request) {
	err := db.URLDB.Ping()
	if err != nil {
		logger.Log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

type (
	RequestBatchItem struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	ResponseBatchItem struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

func postManyURL(w http.ResponseWriter, r *http.Request) {
	var (
		req []RequestBatchItem
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
	longURLs := make([]string, len(req))
	for i, x := range req {
		longURLs[i] = x.OriginalURL
	}
	userID := r.Context().Value(userIDKey{}).(string)
	ids := urlStorage.AddManyURLs(longURLs, userID)
	resp := make([]ResponseBatchItem, len(ids))
	for i, x := range ids {
		resp[i] = ResponseBatchItem{req[i].CorrelationID, config.ResultAddress + "/" + x}
	}

	var body []byte
	if body, err = json.Marshal(resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

type ResponseUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func getUserURLs(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("JWT")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserID(c.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	urls := urlStorage.GetUserURLs(userID)
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var resp = make([]ResponseUserURL, len(urls))
	for i, url := range urls {
		resp[i] = ResponseUserURL{config.ResultAddress + "/" + url[0], url[1]}
	}
	body, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func deleteURLs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey{}).(string)
	var (
		buf bytes.Buffer
		ids []string
	)
	buf.ReadFrom(r.Body)
	defer r.Body.Close()

	json.Unmarshal(buf.Bytes(), &ids)
	urlStorage.DeleteURLs(ids, userID)

	w.WriteHeader(http.StatusAccepted)
}
