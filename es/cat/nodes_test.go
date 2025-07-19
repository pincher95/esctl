package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatNodes(t *testing.T) {
	jsonResp := `[{
        "id":"node-1",
        "ip":"127.0.0.1",
        "heap.percent":"10",
        "ram.percent":"20",
        "cpu":"5",
        "load_1m":"0.1",
        "node.role":"d",
        "node.roles":"data",
        "master":"*",
        "name":"node-1"
    }]`

	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/nodes")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatNodes(ctx, "", "", "", "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(resp) != 1 || resp[0].Name != "node-1" {
		t.Fatalf("unexpected %+v", resp)
	}
}
