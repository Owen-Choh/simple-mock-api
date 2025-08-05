package types

import (
	"encoding/json"
)

type Mapping struct {
	Path     string   `json:"path"`
	Method   string   `json:"method"`
	Response Response `json:"response"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       json.RawMessage   `json:"body,omitempty"`
}
