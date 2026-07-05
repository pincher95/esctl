package node

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestGetNodeStats(t *testing.T) {
	jsonResp := `{"nodes":{"abc123":{
		"name":"es-data-0",
		"jvm":{"mem":{"heap_used_percent":41,"heap_used_in_bytes":1717986918,"heap_max_in_bytes":4294967296},
			"gc":{"collectors":{"young":{"collection_count":120},"old":{"collection_count":3}}}},
		"thread_pool":{"search":{"rejected":5},"write":{"rejected":2}},
		"breakers":{"parent":{"tripped":1},"fielddata":{"tripped":0}},
		"fs":{"total":{"total_in_bytes":250000000000,"available_in_bytes":36000000000}}
	}}}`
	srv, cli := testutil.NewMockServer(jsonResp, "/_nodes/stats/jvm,thread_pool,breaker,fs")
	defer srv.Close()
	shared.SetClient(cli)

	resp, err := GetNodeStats(context.Background(), "")
	if err != nil {
		t.Fatalf("GetNodeStats() error = %v", err)
	}
	n, ok := resp.Nodes["abc123"]
	if !ok {
		t.Fatalf("expected node abc123, got %v", resp.Nodes)
	}
	if n.Name != "es-data-0" {
		t.Errorf("expected name es-data-0, got %q", n.Name)
	}
	if n.JVM.Mem.HeapUsedPercent != 41 {
		t.Errorf("expected heap 41%%, got %d", n.JVM.Mem.HeapUsedPercent)
	}
	if n.ThreadPool["search"].Rejected != 5 {
		t.Errorf("expected 5 search rejections, got %d", n.ThreadPool["search"].Rejected)
	}
	if n.Breakers["parent"].Tripped != 1 {
		t.Errorf("expected parent breaker tripped 1, got %d", n.Breakers["parent"].Tripped)
	}
}
