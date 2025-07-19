package index_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestUpdateIndexSettings(t *testing.T) {
	jsonResp := `{"acknowledged":true}`
	srv, cli := testutil.NewMockServer(jsonResp, "/idx/_settings")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	body := map[string]any{"number_of_replicas": "1"}
	resp, err := idx.UpdateIndexSettings(ctx, "idx", body, true)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if (*resp)["acknowledged"] != true {
		t.Fatalf("unexpected resp %+v", resp)
	}
}
