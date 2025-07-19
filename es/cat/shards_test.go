package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatShards(t *testing.T) {
	jsonResp := `[{
        "index":"idx",
        "shard":"0",
        "prirep":"p",
        "state":"STARTED",
        "docs":"10",
        "store":"1kb",
        "ip":"127.0.0.1",
        "node":"node-1"
    }]`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/shards")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatShards(ctx, "", "", "", "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(resp) != 1 || resp[0].Index != "idx" {
		t.Fatalf("unexpected %+v", resp)
	}
}
