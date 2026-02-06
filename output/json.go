package output

import (
	"encoding/json"
	"fmt"
)

func PrintJson(data any) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to generate pretty JSON: %w", err)
	}

	fmt.Println(string(prettyJSON))
	return nil
}
