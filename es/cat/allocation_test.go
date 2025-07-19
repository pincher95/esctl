package cat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/client"
	"github.com/pincher95/esctl/shared"
)

func TestCatAllocation_BuildsEndpointAndParses(t *testing.T) {
	// Prepare fake response
	jsonResp := `[{
        "shards":"10",
        "disk.indices":"100mb",
        "disk.used":"1gb",
        "disk.avail":"9gb",
        "disk.total":"10gb",
        "disk.percent":"10",
        "host":"node1",
        "ip":"127.0.0.1",
        "node":"node-1"
    }]`

	// Start an httptest server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic path assertion
		if r.URL.Path != "/_cat/allocation" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResp))
	}))
	defer srv.Close()

	// Inject mock client
	cfg := &client.Config{BaseURL: srv.URL}
	shared.SetClient(client.NewClient(cfg))

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatAllocation(ctx, "", "", "")
	if err != nil {
		t.Fatalf("CatAllocation returned error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp))
	}
	if resp[0].Node != "node-1" {
		t.Fatalf("unexpected node: %s", resp[0].Node)
	}
}
