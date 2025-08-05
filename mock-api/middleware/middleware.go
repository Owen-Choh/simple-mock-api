package middleware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request at: %s %s\nRequest: %+v", r.Method, r.URL.Path, r)
		next.ServeHTTP(w, r)
	})
}
