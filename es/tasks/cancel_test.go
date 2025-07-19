package tasks_test

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/es/tasks"
	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestCancelTasks(t *testing.T) {
	jsonResp := `{"nodes":{}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_tasks/_cancel")
	defer srv.Close()
	shared.SetClient(cli)

	ctx := context.Background()
	c := tasks.NewTasks()
	_, err := c.CancelTasks(ctx, "", nil)
	if err != nil {
		t.Fatalf("err %v", err)
	}
}
