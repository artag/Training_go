package main

import (
	"encoding/json" // To convert data to json
	"log"           // To log errors
	"net/http"      // To respond to HTTP requests
	"sync"          // To use the type sync.Mutex
)

const (
	ContentType            = "Content-Type"
	ContentTextPlain       = "text/plain"
	ContentApplicationJson = "application/json"
)

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()
	mu := &sync.Mutex{}

	m.HandleFunc("/", rootHandler)

	t := todoRouter(todoFile, mu)
	m.Handle("/todo", http.StripPrefix("/todo", t))
	m.Handle("/todo/", http.StripPrefix("/todo/", t))

	return m
}

func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set(ContentType, ContentTextPlain)
	w.WriteHeader(status)
	w.Write([]byte(content))
}

func replyJSONContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {
	body, err := json.Marshal(resp)
	if err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set(ContentType, ContentApplicationJson)
	w.WriteHeader(status)
	w.Write(body)
}

func replyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, message)
	http.Error(w, http.StatusText(status), status)
}
