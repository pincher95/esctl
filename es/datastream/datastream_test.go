package datastream

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestList(t *testing.T) {
	mockResp := ListResponse{
		DataStreams: []DataStream{
			{
				Name: "logs-test",
				TimestampField: TimestampField{
					Name: "@timestamp",
				},
				Generation: 1,
				Status:     "GREEN",
				Template:   "logs-template",
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream")
	defer srv.Close()
	shared.SetClient(cli)

	streams, err := List(context.Background(), "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(streams) != 1 {
		t.Errorf("Expected 1 data stream, got %d", len(streams))
	}
	if streams[0].Name != "logs-test" {
		t.Errorf("Expected name 'logs-test', got %s", streams[0].Name)
	}
}

func TestListWithFilter(t *testing.T) {
	mockResp := ListResponse{
		DataStreams: []DataStream{
			{
				Name:       "logs-test",
				Generation: 1,
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream/logs-*")
	defer srv.Close()
	shared.SetClient(cli)

	streams, err := List(context.Background(), "logs-*")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(streams) != 1 {
		t.Errorf("Expected 1 data stream, got %d", len(streams))
	}
}

func TestGet(t *testing.T) {
	mockResp := ListResponse{
		DataStreams: []DataStream{
			{
				Name: "logs-test",
				TimestampField: TimestampField{
					Name: "@timestamp",
				},
				Generation: 3,
				Status:     "GREEN",
				Template:   "logs-template",
				Indices: []DataStreamIndex{
					{
						IndexName: ".ds-logs-test-000001",
						IndexUUID: "uuid1",
					},
				},
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream/logs-test")
	defer srv.Close()
	shared.SetClient(cli)

	stream, err := Get(context.Background(), "logs-test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stream.Name != "logs-test" {
		t.Errorf("Expected name 'logs-test', got %s", stream.Name)
	}
	if stream.Generation != 3 {
		t.Errorf("Expected generation 3, got %d", stream.Generation)
	}
	if len(stream.Indices) != 1 {
		t.Errorf("Expected 1 index, got %d", len(stream.Indices))
	}
}

func TestGetNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_data_stream/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	_, err := Get(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent data stream, got nil")
	}
}

func TestDelete(t *testing.T) {
	mockResp := map[string]interface{}{
		"acknowledged": true,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream/logs-test")
	defer srv.Close()
	shared.SetClient(cli)

	err := Delete(context.Background(), "logs-test")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, "", "/_data_stream/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	err := Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent data stream, got nil")
	}
}

func TestRollover(t *testing.T) {
	mockResp := RolloverResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
		OldIndex:           ".ds-logs-test-000001",
		NewIndex:           ".ds-logs-test-000002",
		RolledOver:         true,
		DryRun:             false,
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream/logs-test/_rollover")
	defer srv.Close()
	shared.SetClient(cli)

	resp, err := Rollover(context.Background(), "logs-test")
	if err != nil {
		t.Fatalf("Rollover() error = %v", err)
	}
	if !resp.RolledOver {
		t.Error("Expected RolledOver to be true")
	}
	if resp.NewIndex != ".ds-logs-test-000002" {
		t.Errorf("Expected new index '.ds-logs-test-000002', got %s", resp.NewIndex)
	}
}

func TestGetStats(t *testing.T) {
	mockResp := StatsResponse{
		DataStreams: []DataStreamStats{
			{
				DataStream:     "logs-test",
				BackingIndices: 3,
				StoreSizeBytes: 1024000,
			},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream/logs-test/_stats")
	defer srv.Close()
	shared.SetClient(cli)

	stats, err := GetStats(context.Background(), "logs-test")
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}
	if len(stats.DataStreams) != 1 {
		t.Errorf("Expected 1 data stream in stats, got %d", len(stats.DataStreams))
	}
	if stats.DataStreams[0].BackingIndices != 3 {
		t.Errorf("Expected 3 backing indices, got %d", stats.DataStreams[0].BackingIndices)
	}
}

func TestGetStatsAll(t *testing.T) {
	mockResp := StatsResponse{
		DataStreams: []DataStreamStats{
			{DataStream: "logs-test1", BackingIndices: 2},
			{DataStream: "logs-test2", BackingIndices: 3},
		},
	}
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_data_stream/_stats")
	defer srv.Close()
	shared.SetClient(cli)

	stats, err := GetStats(context.Background(), "")
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}
	if len(stats.DataStreams) != 2 {
		t.Errorf("Expected 2 data streams in stats, got %d", len(stats.DataStreams))
	}
}
