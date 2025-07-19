package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatHealth(t *testing.T) {
	jsonResp := `[{"status":"green","node.total":"3","shards":"10","relo":"0","init":"0","unassign":"0","pending_tasks":"0","active_shards_percent":"100%"}]`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/health")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatHealth(ctx)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if resp.Status != "green" {
		t.Fatalf("unexpected %+v", resp)
	}
}
