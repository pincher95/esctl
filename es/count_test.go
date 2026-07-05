package es

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCountDocumentsGroupBy(t *testing.T) {
	resp := `{
		"aggregations": {
			"by_index": {
				"buckets": [
					{"key": "test-index", "doc_count": 30, "by_group": {"buckets": [
						{"key": "info", "doc_count": 20},
						{"key": "error", "doc_count": 10}
					]}}
				]
			}
		}
	}`
	srv, cli := testutil.NewMockServer(resp, "/test-index/_search")
	defer srv.Close()
	shared.SetClient(cli)

	counts, err := CountDocuments(context.Background(), "test-index", nil, nil, nil, "level", 0, "", false)
	if err != nil {
		t.Fatalf("CountDocuments() error = %v", err)
	}
	if got := counts["test-index"]["info"]; got != 20 {
		t.Errorf("expected info=20, got %d", got)
	}
	if got := counts["test-index"]["error"]; got != 10 {
		t.Errorf("expected error=10, got %d", got)
	}
}

func TestCountDocumentsTotal(t *testing.T) {
	resp := `{"aggregations": {"by_index": {"buckets": [
		{"key": "test-index", "doc_count": 42}
	]}}}`
	srv, cli := testutil.NewMockServer(resp, "/test-index/_search")
	defer srv.Close()
	shared.SetClient(cli)

	counts, err := CountDocuments(context.Background(), "test-index", nil, nil, nil, "", 0, "", false)
	if err != nil {
		t.Fatalf("CountDocuments() error = %v", err)
	}
	if got := counts["test-index"][""]; got != 42 {
		t.Errorf("expected total=42, got %d", got)
	}
}
