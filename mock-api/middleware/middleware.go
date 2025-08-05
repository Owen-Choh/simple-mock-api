package middleware

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := httputil.DumpRequest(r, true)
		log.Printf("Request received: %s", s)

		next.ServeHTTP(w, r)
	})
}
