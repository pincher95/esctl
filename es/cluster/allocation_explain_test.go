package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestAllocationExplain(t *testing.T) {
	jsonResp := `{"index":"idx","shard":0,"current_state":"UNASSIGNED"}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/allocation/explain")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	resp, err := cluster.ClusterAllocationExplain(ctx, false, false)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if resp.Index != "idx" {
		t.Fatalf("bad resp %+v", resp)
	}
}
