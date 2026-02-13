package index

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestGetIndexStats(t *testing.T) {
	mockResp := IndexStatsResponse{
		Indices: map[string]IndexStat{
			"test-index": {
				UUID: "test-uuid",
				Primaries: IndexStatDetail{
					Docs: IndexStatDocs{
						Count:   1000,
						Deleted: 10,
					},
					Store: IndexStatStore{
						SizeInBytes: 1024000,
					},
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_stats")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetIndexStats(context.Background(), []string{"test-index"}, "")
	if err != nil {
		t.Fatalf("GetIndexStats() error = %v", err)
	}
	if len(resp.Indices) != 1 {
		t.Errorf("Expected 1 index in response, got %d", len(resp.Indices))
	}
}

func TestGetIndexStatsWithMetric(t *testing.T) {
	mockResp := IndexStatsResponse{
		Indices: map[string]IndexStat{
			"test-index": {
				Primaries: IndexStatDetail{
					Docs: IndexStatDocs{
						Count: 1000,
					},
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_stats/docs")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetIndexStats(context.Background(), []string{"test-index"}, "docs")
	if err != nil {
		t.Fatalf("GetIndexStats() error = %v", err)
	}
	if len(resp.Indices) != 1 {
		t.Errorf("Expected 1 index in response, got %d", len(resp.Indices))
	}
}

func TestGetIndexStatsAllIndices(t *testing.T) {
	mockResp := IndexStatsResponse{
		All: IndexStat{
			Primaries: IndexStatDetail{
				Docs: IndexStatDocs{
					Count: 5000,
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_stats")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetIndexStats(context.Background(), []string{}, "")
	if err != nil {
		t.Fatalf("GetIndexStats() error = %v", err)
	}
	if resp.All.Primaries.Docs.Count != 5000 {
		t.Errorf("Expected doc count 5000, got %d", resp.All.Primaries.Docs.Count)
	}
}

func TestGetRecovery(t *testing.T) {
	mockResp := RecoveryResponse{
		"test-index": IndexRecovery{
			Shards: []ShardRecovery{
				{
					ID:    0,
					Type:  "PEER",
					Stage: "DONE",
					Primary: true,
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_recovery")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetRecovery(context.Background(), []string{"test-index"}, false)
	if err != nil {
		t.Fatalf("GetRecovery() error = %v", err)
	}
	if len(resp) != 1 {
		t.Errorf("Expected 1 index in response, got %d", len(resp))
	}
	if len(resp["test-index"].Shards) != 1 {
		t.Errorf("Expected 1 shard, got %d", len(resp["test-index"].Shards))
	}
}

func TestGetRecoveryAllIndices(t *testing.T) {
	mockResp := RecoveryResponse{
		"index1": IndexRecovery{
			Shards: []ShardRecovery{{ID: 0}},
		},
		"index2": IndexRecovery{
			Shards: []ShardRecovery{{ID: 0}},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_recovery")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetRecovery(context.Background(), []string{}, false)
	if err != nil {
		t.Fatalf("GetRecovery() error = %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("Expected 2 indices in response, got %d", len(resp))
	}
}

func TestGetSegments(t *testing.T) {
	mockResp := SegmentsResponse{
		Indices: map[string]IndexSegments{
			"test-index": {
				Shards: map[string][]ShardSegments{
					"0": {
						{
							NumCommittedSegments: 5,
							NumSearchSegments:    5,
						},
					},
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_segments")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetSegments(context.Background(), []string{"test-index"})
	if err != nil {
		t.Fatalf("GetSegments() error = %v", err)
	}
	if len(resp.Indices) != 1 {
		t.Errorf("Expected 1 index in response, got %d", len(resp.Indices))
	}
}

func TestGetSegmentsAllIndices(t *testing.T) {
	mockResp := SegmentsResponse{
		Indices: map[string]IndexSegments{
			"index1": {Shards: map[string][]ShardSegments{}},
			"index2": {Shards: map[string][]ShardSegments{}},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_segments")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetSegments(context.Background(), []string{})
	if err != nil {
		t.Fatalf("GetSegments() error = %v", err)
	}
	if len(resp.Indices) != 2 {
		t.Errorf("Expected 2 indices in response, got %d", len(resp.Indices))
	}
}
