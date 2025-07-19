package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestUpdateClusterSettings(t *testing.T) {
	jsonResp := `{"acknowledged":true}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/settings")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	body := map[string]any{"transient": map[string]any{"indices.recovery.max_bytes_per_sec": "40mb"}}
	resp, err := cluster.UpdateClusterSettings(ctx, body)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if (*resp)["acknowledged"] != true {
		t.Fatalf("unexpected response %+v", resp)
	}
}
