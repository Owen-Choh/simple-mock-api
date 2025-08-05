package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Owen-Choh/simple-mock-api/mock-api/middleware"
	"github.com/Owen-Choh/simple-mock-api/mock-api/utils"
)

type MockServer struct {
	Handler     *ReloadableHandler
	MockPrefix  string
	MappingsDir string
}

func NewMockServer(prefix string, mappingsDir string) *MockServer {
	initialMux := http.NewServeMux()
	return &MockServer{
		Handler:     NewReloadableHandler(initialMux), // Wrap the initial mux in ReloadableHandler
		MockPrefix:  prefix,
		MappingsDir: mappingsDir,
	}
}

// RegisterHandlers loads the initial mappings and sets up the server.
func (ms *MockServer) RegisterHandlers() error {
	newMux, err := ms.createAndLoadMux()
	if err != nil {
		return err
	}

	ms.Handler.UpdateHandler(newMux)

	log.Println("handlers updated.")
	return nil
}

// createAndLoadMux loads mappings from the directory and creates a new http.ServeMux
// with all the mock handlers and the reload handler registered on it.
func (ms *MockServer) createAndLoadMux() (*http.ServeMux, error) {
	newMux := http.NewServeMux()

	mappings, err := utils.LoadMappingsFromDir(ms.MappingsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load mappings from directory '%s': %w", ms.MappingsDir, err)
	}

	for _, mapping := range mappings {
		pattern := mapping.Method + " " + ms.MockPrefix + mapping.Path
		newMux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			utils.WriteResponse(w, mapping.Response)
		})
	}

	newMux.HandleFunc("/__admin/reload", ms.handleReload)

	log.Printf("created a new ServeMux with %d mock mappings.", len(mappings))
	return newMux, nil
}

// handleReload is the HTTP handler for the /__admin/reload endpoint.
// It triggers the reloading of mappings and updates the active handler.
func (ms *MockServer) handleReload(w http.ResponseWriter, r *http.Request) {
	log.Println("received reload request. attempting to reload mappings...")
	
	if err := ms.RegisterHandlers(); err != nil {
		http.Error(w, fmt.Sprintf("Error reloading mappings: %s", err.Error()), http.StatusInternalServerError)
		log.Printf("failed to reload mappings: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "mappings reloaded successfully")
	log.Println("mappings reloaded successfully.")
}

// Start begins the HTTP server listening on the specified address.
func (ms *MockServer) Start(addr string) error {
	if addr == "" {
		return fmt.Errorf("address cannot be empty")
	}

	handlerWithMiddleware := middleware.LoggingMiddleware(ms.Handler)

	server := &http.Server{
		Addr:    addr,
		Handler: handlerWithMiddleware,
	}

	log.Printf("starting mock server on %s...\n", addr)
	return server.ListenAndServe()
}
