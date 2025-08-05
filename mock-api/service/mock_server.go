package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/simple-mock-api/mock-api/middleware"
	"github.com/Owen-Choh/simple-mock-api/mock-api/types"
	"github.com/Owen-Choh/simple-mock-api/mock-api/utils"
)

type MockServer struct {
	Server *http.ServeMux
}

func NewMockServer() *MockServer {
	return &MockServer{
		Server: http.NewServeMux(),
	}
}

func (ms *MockServer) LoadMappings(prefix string, mappings []types.Mapping) error {
	if len(mappings) == 0 {
		return fmt.Errorf("no mappings to load")
	}

	for _, mapping := range mappings {
		pattern := mapping.Method + " " + prefix + mapping.Path
		ms.Server.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Handling request for: %s %s", r.Method, r.URL.Path)
			utils.WriteResponse(w, mapping.Response)
		})
	}
	return nil
}

func (ms *MockServer) Start(addr string) error {
	if addr == "" {
		return fmt.Errorf("address cannot be empty")
	}

	handlerWithMiddleware := middleware.LoggingMiddleware(ms.Server)

	server := &http.Server{
		Addr:    addr,
		Handler: handlerWithMiddleware,
	}

	log.Printf("Starting mock server on %s...\n", addr)
	return server.ListenAndServe()
}
