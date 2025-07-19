package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatFielddata(t *testing.T) {
	jsonResp := `[{
        "id":"node-1",
        "host":"host",
        "ip":"127.0.0.1",
        "node":"node1",
        "field":"title",
        "size":"1kb"
    }]`

	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/fielddata")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatFielddata(ctx, "", "", "")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(resp) != 1 || resp[0].Field != "title" {
		t.Fatalf("unexpected response %+v", resp)
	}
}
