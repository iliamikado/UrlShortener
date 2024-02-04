package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodPOST(t *testing.T) {
	urlsMap = make(map[string]string)

	t.Run("right request", func(t *testing.T) {
		longURL := "https://ya.ru"
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL))
		w := httptest.NewRecorder()
		mainPage(w, req)
		assert.Equal(t, http.StatusCreated, w.Code, "Wrong status code")
	
		shortURL := w.Body.String()
		assert.NotNil(t, shortURL, "No short URL in response")
		assert.Contains(t, shortURL, "http://example.com/", "Short URL should contains server adress, got " + shortURL)
	})

	t.Run("without body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		mainPage(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code, "Wrong status code")
	})
}

func TestMethodGET(t *testing.T) {
	urlsMap = make(map[string]string)

	t.Run("right request", func(t *testing.T) {
		longURL := "https://ya.ru"
		postReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL))
		postW := httptest.NewRecorder()
		mainPage(postW, postReq)
		shortURL := postW.Body.String()
		getReq := httptest.NewRequest(http.MethodGet, shortURL, nil)
		getW := httptest.NewRecorder()
		mainPage(getW, getReq)
		assert.Equal(t, http.StatusTemporaryRedirect, getW.Code, "Wrong status code")
		assert.Equal(t, longURL, getW.Header().Get("Location"), "Wrong long url")
	})

	t.Run("request unexisted url", func(t *testing.T) {
		getReq := httptest.NewRequest(http.MethodGet, "/123", nil)
		getW := httptest.NewRecorder()
		mainPage(getW, getReq)
		assert.Equal(t, http.StatusBadRequest, getW.Code)
	})
}


