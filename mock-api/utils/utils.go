package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Owen-Choh/simple-mock-api/mock-api/types"
	"net/http"
	"os"
)

func ParseJsonFile(jsonFile *os.File, v interface{}) error {
	jsonParser := json.NewDecoder(jsonFile)
	if err := jsonParser.Decode(v); err != nil {
		return fmt.Errorf("parsing json file: %w", err)
	}
	return nil
}

func WriteResponse(w http.ResponseWriter, r types.Response) error {
	w.WriteHeader(r.StatusCode)

	if r.Headers != nil {
		for key, value := range r.Headers {
			w.Header().Set(key, value)
		}
	}

	if r.Body != nil {
		_, err := w.Write(r.Body)
		return err
	}

	return nil
}
