package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJsonFile(jsonFile *os.File, v interface{}) error {
	jsonParser := json.NewDecoder(jsonFile)
	if err := jsonParser.Decode(v); err != nil {
		return fmt.Errorf("parsing json file: %w", err)
	}
	return nil
}
