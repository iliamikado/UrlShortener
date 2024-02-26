package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iliamikado/UrlShortener/internal/config"
	"github.com/iliamikado/UrlShortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMethodPOST(t *testing.T) {
	srv := launchServer()
	defer srv.Close()

	t.Run("right request", func(t *testing.T) {
		longURL := "https://ya.ru"
		resp, shortURL := testRequest(t, srv, http.MethodPost, "/", longURL)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Wrong status code")
		assert.NotNil(t, shortURL, "No short URL in response")
		assert.Contains(t, shortURL, srv.URL, "Short URL should contains server adress, got " + shortURL)
	})

	t.Run("without body", func(t *testing.T) {
		resp, _ := testRequest(t, srv, http.MethodPost, "/", "")
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Wrong status code")
	})
}

func TestMethodGET(t *testing.T) {
	srv := launchServer()
	defer srv.Close()

	t.Run("right request", func(t *testing.T) {
		longURL := "https://ya.ru"
		postResp, shortURL := testRequest(t, srv, http.MethodPost, "/", longURL)
		defer postResp.Body.Close()
		shortURL = strings.Replace(shortURL, srv.URL, "", 1)
		resp, _ := testRequest(t, srv, http.MethodGet, shortURL, "")
		defer resp.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Wrong status code")
		assert.Equal(t, longURL, resp.Header.Get("Location"), "Wrong long url")
	})

	t.Run("request unexisted url", func(t *testing.T) {
		resp, _ := testRequest(t, srv, http.MethodGet, "/123", "")
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func launchServer() *httptest.Server {
	urlStorage = storage.NewURLStorage()
	srv := httptest.NewServer(AppRouter())
	config.LaunchAddress = srv.URL
	config.ResultAddress = srv.URL
	return srv
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

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
