package utils

import (
	"encoding/json"
	"fmt"
)

// FormatSettingValue renders settings values for CLI table output.
func FormatSettingValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		if raw, err := json.Marshal(v); err == nil {
			return string(raw)
		}
		return fmt.Sprint(v)
	}
}
