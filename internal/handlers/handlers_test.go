package handlers

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iliamikado/UrlShortener/internal/storage"
)

func CreateURLStorage() {
	urlStorage = storage.NewSimpleStorage()
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func randomURL() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "https://" + string(b) + ".com"
}

func CreateReqAndRes(method, path string, body io.Reader) (http.ResponseWriter, *http.Request) {
	r := httptest.NewRequest(http.MethodPost, path, body)
	r = r.WithContext(context.WithValue(r.Context(), userIDKey{}, "default"))
	w := httptest.NewRecorder()
	return w, r
}

func BenchmarkPostURL(b *testing.B) {
	CreateURLStorage()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		w, r := CreateReqAndRes(http.MethodPost, "/", strings.NewReader(randomURL()))
		b.StartTimer()
		PostURL(w, r)
	}
}

func BenchmarkGetShortURL(b *testing.B) {
	CreateURLStorage()
	userID := "default"
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		shortURL, _ := urlStorage.AddURL(randomURL(), userID)
		w, r := CreateReqAndRes(http.MethodGet, "/"+shortURL, nil)
		b.StartTimer()
		GetURL(w, r)
	}
}
