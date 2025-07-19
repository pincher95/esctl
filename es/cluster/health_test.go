package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestClusterHealth(t *testing.T) {
	jsonResp := `{"status":"green","number_of_nodes":1}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/health")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	h, err := cluster.ClusterHealth(ctx, "", "", "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if h.Status != "green" {
		t.Fatalf("bad status %s", h.Status)
	}
}
