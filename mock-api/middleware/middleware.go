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

type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// DefaultCORSOptions allows everything
var DefaultCORSOptions = CORSOptions{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"Content-Type", "Authorization"},
}

// CORS middleware
func CORS(opts CORSOptions) func(http.Handler) http.Handler {
	originMap := make(map[string]bool)
	for _, o := range opts.AllowedOrigins {
		originMap[o] = true
	}

	allowMethods := join(opts.AllowedMethods)
	allowHeaders := join(opts.AllowedHeaders)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Allow any origin if wildcard, otherwise match explicitly
			if originMap["*"] || originMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func join(list []string) string {
	result := ""
	for i, s := range list {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
