package output

import (
	"testing"

	"github.com/pincher95/esctl/shared"
)

func TestRenderUnsupportedFormat(t *testing.T) {
	previous := shared.OutputFormat
	t.Cleanup(func() { shared.OutputFormat = previous })

	shared.OutputFormat = "xml"
	if err := Render(map[string]string{"a": "b"}); err == nil {
		t.Fatalf("expected error for unsupported output format")
	}
}

func TestRenderTableDefaultsToJSON(t *testing.T) {
	previous := shared.OutputFormat
	t.Cleanup(func() { shared.OutputFormat = previous })

	shared.OutputFormat = "table"
	if err := Render(map[string]string{"a": "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintJsonError(t *testing.T) {
	if err := PrintJson(func() {}); err == nil {
		t.Fatalf("expected JSON marshal error")
	}
}

func TestPrintYamlError(t *testing.T) {
	if err := PrintYaml(func() {}); err == nil {
		t.Fatalf("expected YAML marshal error")
	}
}
