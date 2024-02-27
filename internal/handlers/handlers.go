package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/iliamikado/UrlShortener/internal/config"
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
	id := urlStorage.AddURL(req.URL)
	shortURL := config.ResultAddress + "/" + id
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

func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ow := w

        acceptEncoding := r.Header.Get("Accept-Encoding")
        supportsGzip := strings.Contains(acceptEncoding, "gzip")
        if supportsGzip {
            cw := newCompressWriter(w)
            ow = cw
            defer cw.Close()
        }

        contentEncoding := r.Header.Get("Content-Encoding")
        sendsGzip := strings.Contains(contentEncoding, "gzip")
        if sendsGzip {
            cr, err := newCompressReader(r.Body)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            r.Body = cr
            defer cr.Close()
        }
        next.ServeHTTP(ow, r)
    })
}
