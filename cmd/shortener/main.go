package main

import (
	"io"
	"net/http"
	"github.com/iliamikado/UrlShortener/internal/logic"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

var urlsMap map[string]string
func run() error {
	urlsMap = make(map[string]string)
	mx := http.NewServeMux()
	mx.HandleFunc("/", mainPage)
	return http.ListenAndServe(":8080", mx)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postURL(w, r)
	} else if r.Method == http.MethodGet {
		getURL(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest);
	}
}

func postURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	longURL := string(body)
	id := logic.AddURL(urlsMap, longURL)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	shortUrl := scheme + "://" + r.Host + "/" + id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func getURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]
	longURL, err := logic.GetURL(urlsMap, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}


