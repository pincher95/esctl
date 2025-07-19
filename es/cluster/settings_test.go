package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestClusterSettings(t *testing.T) {
	jsonResp := `{"persistent":{},"transient":{}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/settings")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	s, err := cluster.ClusterSettings(ctx, true, false)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if s == nil {
		t.Fatalf("settings nil")
	}
}
