package handlers

import (
	"net/http"
	"strings"
)

func ExamplePostURL() {
	CreateURLStorage()
	url := "http://ya.ru"
	w, r := CreateReqAndRes(http.MethodPost, "/", strings.NewReader(url))
	PostURL(w, r)
}
