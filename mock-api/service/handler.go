package service

import (
	"net/http"
	"sync"
)

// Custom http.Handler that can be dynamically updated.
type ReloadableHandler struct {
	mu             sync.RWMutex // Mutex to protect currentHandler during reads and writes
	currentHandler http.Handler
}

func NewReloadableHandler(initialHandler http.Handler) *ReloadableHandler {
	return &ReloadableHandler{
		currentHandler: initialHandler,
	}
}

func (rh *ReloadableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.mu.RLock()
	handler := rh.currentHandler
	rh.mu.RUnlock()

	handler.ServeHTTP(w, r)
}

func (rh *ReloadableHandler) UpdateHandler(newHandler http.Handler) {
	rh.mu.Lock()
	rh.currentHandler = newHandler
	rh.mu.Unlock()
}