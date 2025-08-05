package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func LoadMappingsFromFile(filePath string) ([]types.Mapping, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening mappings file: %w", err)
	}
	defer file.Close()

	var mappings []types.Mapping
	if err := ParseJsonFile(file, &mappings); err != nil {
		return nil, fmt.Errorf("parsing mappings file: %w", err)
	}

	return mappings, nil
}

func LoadMappingsFromDir(dir string) ([]types.Mapping, error) {
	var allMappings []types.Mapping

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking file path: %w", err)
		}

		// Skip directories and non-json files
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		var mappings []types.Mapping
		mappings, err = LoadMappingsFromFile(path)
		if err != nil {
			return fmt.Errorf("loading mappings from file %s: %w", path, err)
		}

		allMappings = append(allMappings, mappings...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allMappings, nil
}

// compares two slices of Mapping for equality.
func MappingsEqual(a, b []types.Mapping) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Path != b[i].Path || !strings.EqualFold(a[i].Method, b[i].Method) ||
			a[i].Response.StatusCode != b[i].Response.StatusCode {
			return false
		}

		if !bytes.Equal(
			canonicalizeJSON(a[i].Response.Body),
			canonicalizeJSON(b[i].Response.Body),
		) {
			return false
		}
	}

	return true
}

// removes insignificant whitespace from JSON for consistent comparison.
func canonicalizeJSON(input json.RawMessage) json.RawMessage {
	var compacted bytes.Buffer
	json.Compact(&compacted, input)
	return compacted.Bytes()
}
