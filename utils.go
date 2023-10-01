package main

import (
	"log"
	"net/http"
)

func handleLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
