package index_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestGetAliases(t *testing.T) {
	jsonResp := `{"idx":{"aliases":{"alias1":{}}}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_alias/_all")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	resp, err := idx.GetAliases(ctx, "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if _, ok := (*resp)["idx"]; !ok {
		t.Fatalf("alias missing in response")
	}
}
