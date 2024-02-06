package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethodPOST(t *testing.T) {
	urlsMap = make(map[string]string)
	srv := httptest.NewServer(AppRouter())
	defer srv.Close()

	t.Run("right request", func(t *testing.T) {
		longURL := "https://ya.ru"
		resp, shortURL := testRequest(t, srv, http.MethodPost, "/", longURL)
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Wrong status code")
		assert.NotNil(t, shortURL, "No short URL in response")
		assert.Contains(t, shortURL, srv.URL, "Short URL should contains server adress, got " + shortURL)
	})

	t.Run("without body", func(t *testing.T) {
		resp, _ := testRequest(t, srv, http.MethodPost, "/", "")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Wrong status code")
	})
}

func TestMethodGET(t *testing.T) {
	urlsMap = make(map[string]string)
	srv := httptest.NewServer(AppRouter())
	defer srv.Close()

	t.Run("right request", func(t *testing.T) {
		longURL := "https://ya.ru"
		_, shortURL := testRequest(t, srv, http.MethodPost, "/", longURL)
		shortURL = strings.Replace(shortURL, srv.URL, "", 1)
		resp, _ := testRequest(t, srv, http.MethodGet, shortURL, "")
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Wrong status code")
		assert.Equal(t, longURL, resp.Header.Get("Location"), "Wrong long url")
	})

	t.Run("request unexisted url", func(t *testing.T) {
		resp, _ := testRequest(t, srv, http.MethodGet, "/123", "")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func testRequest(t *testing.T, ts *httptest.Server, method, path, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
