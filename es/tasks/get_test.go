package tasks_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/tasks"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestGetTasks(t *testing.T) {
	jsonResp := `{"nodes":{"node-1":{"name":"node-1","tasks":{"123":{"node":"node-1","id":123,"action":"search","description":"desc"}}}}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_tasks")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := tasks.NewTasks()
	resp, err := c.GetTasks(ctx, "", nil)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(resp.Nodes) != 1 {
		t.Fatalf("unexpected resp %+v", resp)
	}
}
