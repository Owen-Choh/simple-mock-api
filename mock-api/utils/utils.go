package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Owen-Choh/simple-mock-api/mock-api/types"
)

func ParseJsonFile(jsonFile *os.File, v interface{}) error {
	jsonParser := json.NewDecoder(jsonFile)
	if err := jsonParser.Decode(v); err != nil {
		return fmt.Errorf("parsing json file: %w", err)
	}
	return nil
}

func WriteResponse(w http.ResponseWriter, r types.Response) error {
	for key, value := range r.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(r.StatusCode)

	return json.NewEncoder(w).Encode(r.Body)
}
