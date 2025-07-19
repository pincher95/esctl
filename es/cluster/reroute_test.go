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
