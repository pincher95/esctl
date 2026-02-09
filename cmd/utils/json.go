package utils

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON decodes JSON and returns a consistent error prefix.
func UnmarshalJSON(data []byte, target any, errorPrefix string) error {
	if err := json.Unmarshal(data, target); err != nil {
		if errorPrefix == "" {
			return fmt.Errorf("invalid JSON: %w", err)
		}
		return fmt.Errorf("%s: %w", errorPrefix, err)
	}
	return nil
}

// ParseJSONMap parses JSON into a map with a consistent error prefix.
func ParseJSONMap(input string, errorPrefix string) (map[string]any, error) {
	var out map[string]any
	if err := UnmarshalJSON([]byte(input), &out, errorPrefix); err != nil {
		return nil, err
	}
	return out, nil
}
