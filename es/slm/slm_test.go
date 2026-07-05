package slm

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestListSLMPolicies(t *testing.T) {
	jsonResp := `{
		"daily-snapshots": {
			"version": 1,
			"modified_date": "2024-01-01T00:00:00Z",
			"next_execution": "2024-01-02T01:30:00Z",
			"policy": {"name": "<daily-{now/d}>", "schedule": "0 30 1 * * ?", "repository": "my_repo",
				"config": {"indices": ["*"]}, "retention": {"expire_after": "30d"}}
		}
	}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_slm/policy")
	defer srv.Close()
	shared.SetClient(cli)

	policies, err := List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	p, ok := policies["daily-snapshots"]
	if !ok {
		t.Fatalf("expected daily-snapshots, got %v", policies)
	}
	if p.Name != "daily-snapshots" {
		t.Errorf("expected Name set from key, got %q", p.Name)
	}
	if p.Policy.Repository != "my_repo" || p.Policy.Schedule != "0 30 1 * * ?" {
		t.Errorf("unexpected policy definition: %+v", p.Policy)
	}
}

func TestExecuteSLMPolicy(t *testing.T) {
	srv, cli := testutil.NewMockServer(`{"snapshot_name":"daily-snap-20240101"}`, "/_slm/policy/daily-snapshots/_execute")
	defer srv.Close()
	shared.SetClient(cli)

	resp, err := Execute(context.Background(), "daily-snapshots")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resp.SnapshotName != "daily-snap-20240101" {
		t.Errorf("expected snapshot name, got %q", resp.SnapshotName)
	}
}
