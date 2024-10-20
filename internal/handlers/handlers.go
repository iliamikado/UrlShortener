// Пакет, где реализованы хэндлеры и миддлвары
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

// Создает Router с хендлерами, нужными для приложения и использующий переданное хранилище.
func AppRouter(st storage.URLStorage) *chi.Mux {
	urlStorage = st

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(gzipMiddleware)
	r.Get("/{id}", GetURL)
	r.Post("/", authMiddleware(PostURL))
	r.Post("/api/shorten", authMiddleware(PostJSON))
	r.Get("/ping", PingDB)
	r.Post("/api/shorten/batch", authMiddleware(PostManyURL))
	r.Get("/api/user/urls", GetUserURLs)
	r.Delete("/api/user/urls", authMiddleware(DeleteURLs))

	return r
}

func createShortURL(longURL string, userID string) (string, error) {
	logger.Log.Info("Get longUrl = " + longURL)
	id, err := urlStorage.AddURL(longURL, userID)
	shortURL := config.ResultAddress + "/" + id
	return shortURL, err
}

// postURL сохраняет длинный url и возвращает короткий.
func PostURL(w http.ResponseWriter, r *http.Request) {
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

// getURL обрабатывает короткие url и возвращает длинные.
func GetURL(w http.ResponseWriter, r *http.Request) {
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

// Структуры для json запросов
type (
	RequestJSON struct {
		URL string `json:"url"`
	}
	ResponseJSON struct {
		Result string `json:"result"`
	}
)

// postJSON сохраняет длинный url и возвращает короткий. Получение и отправка данных происходит в формате json.
func PostJSON(w http.ResponseWriter, r *http.Request) {
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

// pingDB проверяет есть ли соединение с базой данных.
func PingDB(w http.ResponseWriter, r *http.Request) {
	err := db.URLDB.Ping()
	if err != nil {
		logger.Log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Структуры для множественного запроса
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

// postManyURL сокращает одновременно несколько длинных url. Получение и отправка в формате json.
func PostManyURL(w http.ResponseWriter, r *http.Request) {
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

// ResponseUserURL - ответ пользователю
type ResponseUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// getUserURLs возвращает все urls (короткие и длинные версии) сохранненые пользователем.
func GetUserURLs(w http.ResponseWriter, r *http.Request) {
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

// deleteURLs удаляет urls созданные пользователем с переданными id.
func DeleteURLs(w http.ResponseWriter, r *http.Request) {
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
