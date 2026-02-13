package index

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestOpen(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_open")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.Open(context.Background(), []string{"test-index"})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
}

func TestOpenMultipleIndices(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/index1,index2/_open")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.Open(context.Background(), []string{"index1", "index2"})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
}

func TestOpenNoIndices(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Open(context.Background(), []string{})
	if err == nil {
		t.Error("Expected error when no indices specified, got nil")
	}
}

func TestClose(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/test-index/_close")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.Close(context.Background(), []string{"test-index"})
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
}

func TestCloseNoIndices(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Close(context.Background(), []string{})
	if err == nil {
		t.Error("Expected error when no indices specified, got nil")
	}
}

func TestClone(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
		Index:              "target-index",
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/source-index/_clone/target-index")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.Clone(context.Background(), "source-index", "target-index", nil)
	if err != nil {
		t.Fatalf("Clone() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
	if resp.Index != "target-index" {
		t.Errorf("Expected Index = target-index, got %s", resp.Index)
	}
}

func TestCloneWithSettings(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
		Index:              "target-index",
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/source-index/_clone/target-index")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	settings := map[string]any{
		"index.number_of_replicas": 2,
	}
	resp, err := idx.Clone(context.Background(), "source-index", "target-index", settings)
	if err != nil {
		t.Fatalf("Clone() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
}

func TestCloneNoSource(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Clone(context.Background(), "", "target", nil)
	if err == nil {
		t.Error("Expected error when source not specified, got nil")
	}
}

func TestCloneNoTarget(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Clone(context.Background(), "source", "", nil)
	if err == nil {
		t.Error("Expected error when target not specified, got nil")
	}
}

func TestSplit(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
		Index:              "target-index",
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/source-index/_split/target-index")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.Split(context.Background(), "source-index", "target-index", 6, nil)
	if err != nil {
		t.Fatalf("Split() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
}

func TestSplitInvalidShards(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Split(context.Background(), "source", "target", 0, nil)
	if err == nil {
		t.Error("Expected error when shards <= 0, got nil")
	}
}

func TestShrink(t *testing.T) {
	mockResp := IndexOperationResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
		Index:              "target-index",
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/source-index/_shrink/target-index")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.Shrink(context.Background(), "source-index", "target-index", 1, nil)
	if err != nil {
		t.Fatalf("Shrink() error = %v", err)
	}
	if !resp.Acknowledged {
		t.Errorf("Expected Acknowledged = true, got %v", resp.Acknowledged)
	}
}

func TestShrinkInvalidShards(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Shrink(context.Background(), "source", "target", -1, nil)
	if err == nil {
		t.Error("Expected error when shards <= 0, got nil")
	}
}

func TestShrinkNoSource(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Shrink(context.Background(), "", "target", 1, nil)
	if err == nil {
		t.Error("Expected error when source not specified, got nil")
	}
}

func TestShrinkNoTarget(t *testing.T) {
	idx := NewIndex()
	_, err := idx.Shrink(context.Background(), "source", "", 1, nil)
	if err == nil {
		t.Error("Expected error when target not specified, got nil")
	}
}
