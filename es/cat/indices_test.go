package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatIndices(t *testing.T) {
	jsonResp := `[{
        "health":"green",
        "status":"open",
        "index":"idx",
        "uuid":"uuid1",
        "pri":"1",
        "rep":"1",
        "docs.count":"10",
        "docs.deleted":"0",
        "creation.date.string":"now",
        "store.size":"1kb",
        "pri.store.size":"500b"
    }]`

	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/indices")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatIndices(ctx, "", "", "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(resp) != 1 || resp[0].Index != "idx" {
		t.Fatalf("bad resp %+v", resp)
	}
}
