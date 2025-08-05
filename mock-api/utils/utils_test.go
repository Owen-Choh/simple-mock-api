package utils_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Owen-Choh/simple-mock-api/mock-api/types"
	"github.com/Owen-Choh/simple-mock-api/mock-api/utils"
)

func TestWriteResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	resp := types.Response{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: json.RawMessage(`{"message": "created"}`),
	}

	err := utils.WriteResponse(rec, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Errorf("expected status 201, got %d", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	if body["message"] != "created" {
		t.Errorf("unexpected body: %v", body)
	}
}

func TestParseJsonFile(t *testing.T) {
	tempDirectory := t.TempDir()

	testCases := []struct {
		name  string
		input []types.Mapping
	}{
		{
			name: "GET mapping with JSON body",
			input: []types.Mapping{
				{
					Path:   "/test",
					Method: "GET",
					Response: types.Response{
						StatusCode: 200,
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						Body: json.RawMessage(`{"message":"ok"}`),
					},
				},
			},
		},
		{
			name: "GET mapping with string body",
			input: []types.Mapping{
				{
					Path:   "/test",
					Method: "GET",
					Response: types.Response{
						StatusCode: 300,
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						Body: json.RawMessage(`"ok"`),
					},
				},
			},
		},
		{
			name: "POST mapping with string body",
			input: []types.Mapping{
				{
					Path:   "/login",
					Method: "POST",
					Response: types.Response{
						StatusCode: 200,
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						Body: json.RawMessage(`"ok"`),
					},
				},
			},
		},
		{
			name: "mapping without header",
			input: []types.Mapping{
				{
					Path:   "/test/no-header",
					Method: "GET",
					Response: types.Response{
						StatusCode: 200,
						Body:       json.RawMessage(`"ok"`),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile := createTempFileWithMapping(t, tempDirectory, "testfile_*.json", tc.input)
			defer tmpFile.Close()

			var loaded []types.Mapping
			err := utils.ParseJsonFile(tmpFile, &loaded)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !utils.MappingsEqual(loaded, tc.input) {
				t.Fatalf("unexpected mapping: %+v", loaded)
			}
		})
	}
}

func TestLoadMappingsFromFile(t *testing.T) {
	tempDirectory := t.TempDir()

	mappings := []types.Mapping{
		{
			Path:   "/hello",
			Method: "GET",
			Response: types.Response{
				StatusCode: 201,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: json.RawMessage(`"hello"`),
			},
		},
		{
			Path:   "/login/",
			Method: "POST",
			Response: types.Response{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: json.RawMessage(`{"message": "hello"}`),
			},
		},
	}

	tmpFile := createTempFileWithMapping(t, tempDirectory, "mapping_*.json", mappings)
	defer tmpFile.Close()

	loaded, err := utils.LoadMappingsFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadMappingsFromFile failed: %v", err)
	}

	if !utils.MappingsEqual(loaded, mappings) {
		t.Errorf("expected: %+v, got: %+v", mappings, loaded)
	}
}

func TestLoadMappingsFromDir(t *testing.T) {
	tempDirectory := t.TempDir()

	mapping1 := []types.Mapping{
		{
			Path:   "/one",
			Method: "GET",
			Response: types.Response{
				StatusCode: 200,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: json.RawMessage(`{"msg":"one"}`),
			},
		},
	}
	mapping2 := []types.Mapping{
		{
			Path:   "/two",
			Method: "POST",
			Response: types.Response{
				StatusCode: 201,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: json.RawMessage(`{"msg":"two"}`),
			},
		},
	}
	mapping_file_1 := createTempFileWithMapping(t, tempDirectory, "mapping_one_*.json", mapping1)
	mapping_file_2 := createTempFileWithMapping(t, tempDirectory, "mapping_two_*.json", mapping2)
	skipped_file := createTempFile(t, tempDirectory, "skip_*.txt", "not json") // should be ignored as it is not a JSON file
	defer mapping_file_1.Close()
	defer mapping_file_2.Close()
	defer skipped_file.Close()

	full_mappings := append(mapping1, mapping2...)
	loaded, err := utils.LoadMappingsFromDir(tempDirectory)
	if err != nil {
		t.Fatalf("error loading mappings from dir: %v", err)
	}

	if !utils.MappingsEqual(loaded, full_mappings) {
		t.Errorf("%s expected: %+v, got: %+v",tempDirectory, full_mappings, loaded)
	}
}

// --- Helpers ---

// creates a temporary file with the specified mappings. the created file is expected to be closed by the caller.
func createTempFileWithMapping(t *testing.T, dir string, filename string, mappings []types.Mapping) *os.File {
	// Marshal to JSON to write to file
	jsonBytes, err := json.Marshal(mappings)
	if err != nil {
		t.Fatalf("failed to marshal mappings: %v", err)
	}

	// Create temp file with the JSON data
	return createTempFile(t, dir, filename, string(jsonBytes))
}

// creates a temporary file with the specified name and content. the created file is expected to be closed by the caller.
func createTempFile(t *testing.T, dir string, name string, content string) *os.File {
	tmpFile, err := os.CreateTemp(dir, name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek to beginning of temp file: %v", err)
	}

	return tmpFile
}
