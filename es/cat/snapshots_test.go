package cat_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCatSnapshots(t *testing.T) {
	jsonResp := `[{
        "id":"snap-1",
        "status":"SUCCESS",
        "start_epoch":"1",
        "start_time":"now",
        "end_epoch":"2",
        "end_time":"now",
        "duration":"1s",
        "indices":"1",
        "successful_shards":"1",
        "failed_shards":"0",
        "total_shards":"1",
        "reason":""
    }]`

	srv, cli := testutil.NewMockServer(jsonResp, "/_cat/snapshots")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := cat.NewCat()
	resp, err := c.CatSnapshots(ctx, "", "", "")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(resp) != 1 || resp[0].ID != "snap-1" {
		t.Fatalf("unexpected %+v", resp)
	}
}
