package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatPlugins(t *testing.T) {
	jsonResp := `[{
        "id":"node-1",
        "name":"analysis-icu",
        "component":"analysis-icu",
        "version":"1.0",
        "description":"icu"
    }]`
	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/plugins")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatPlugins(ctx, "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(resp) != 1 || resp[0].Name != "analysis-icu" {
		t.Fatalf("unexpected %+v", resp)
	}
}
