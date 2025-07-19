package index_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCacheClear(t *testing.T) {
	jsonResp := `{"_shards":{"total":1,"successful":1,"failed":0}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cache/clear")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	resp, err := idx.CacheClear(ctx, "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if resp.Shards.Successful != 1 {
		t.Fatalf("unexpected shards %+v", resp.Shards)
	}
}
