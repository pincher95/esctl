package output

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/shared"
)

// Render prints data in the configured output format.
func Render(data any) error {
	switch strings.ToLower(shared.OutputFormat) {
	case "", "table", "json":
		return PrintJson(data)
	case "yaml":
		return PrintYaml(data)
	default:
		return fmt.Errorf("unsupported output format: %s", shared.OutputFormat)
	}
}
