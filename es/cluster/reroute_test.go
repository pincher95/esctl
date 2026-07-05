package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestClusterReroute(t *testing.T) {
	jsonResp := `{"acknowledged":true}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/reroute")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	resp, err := cluster.ClusterReroute(ctx, "", false, false, false)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if !resp.Acknowledged {
		t.Fatalf("expected acknowledged")
	}
}

func TestClusterRerouteCommands(t *testing.T) {
	jsonResp := `{"acknowledged":true}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/reroute")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	commands := []map[string]any{
		{"allocate_stale_primary": map[string]any{"index": "my-index", "shard": 0, "node": "es-data-1", "accept_data_loss": true}},
	}
	resp, err := cluster.ClusterRerouteCommands(ctx, commands, true, true)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if !resp.Acknowledged {
		t.Fatalf("expected acknowledged")
	}
}
