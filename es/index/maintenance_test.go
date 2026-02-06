package index_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

const shardsJSON = `{"_shards":{"total":5,"successful":5,"failed":0}}`

func TestRefresh_All(t *testing.T) {
	srv, cli := testutil.NewMockServer(shardsJSON, "/_refresh")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	resp, err := idx.Refresh(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Shards.Successful != 5 {
		t.Fatalf("unexpected shards: %+v", resp.Shards)
	}
}

func TestRefresh_Index(t *testing.T) {
	srv, cli := testutil.NewMockServer(shardsJSON, "/my-index/_refresh")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	_, err := idx.Refresh(ctx, "my-index")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFlush_All(t *testing.T) {
	srv, cli := testutil.NewMockServer(shardsJSON, "/_flush")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	_, err := idx.Flush(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForcemerge_Index(t *testing.T) {
	srv, cli := testutil.NewMockServer(shardsJSON, "/idx/_forcemerge")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	idx := index.NewIndex()
	_, err := idx.Forcemerge(ctx, "idx", 1, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
