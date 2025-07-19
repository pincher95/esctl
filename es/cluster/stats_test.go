package cluster_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestClusterStats(t *testing.T) {
	jsonResp := `{"cluster_name":"test","status":"green","indices":{"count":1},"_nodes":{"total":1,"successful":1,"failed":0}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cluster/stats")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	s, err := cluster.ClusterStats(ctx, "", false)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if s.ClusterName != "test" {
		t.Fatalf("bad name %s", s.ClusterName)
	}
}
