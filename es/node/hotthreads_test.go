package node

import (
	"context"
	"strings"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestGetHotThreads(t *testing.T) {
	mockResponse := `::: {node-1}{abc123}{127.0.0.1}{127.0.0.1:9300}
   Hot threads at 2024-01-01T00:00:00.000Z, interval=500ms, busiestThreads=3, ignoreIdleThreads=true:

   99.9% (499.5ms out of 500ms) cpu usage by thread 'elasticsearch[node-1][search][T#1]'
     10/10 snapshots sharing following 5 elements
       java.lang.Thread.run(Thread.java:750)

   0.1% (0.5ms out of 500ms) cpu usage by thread 'elasticsearch[node-1][write][T#2]'
     2/10 snapshots sharing following 3 elements
       java.lang.Thread.run(Thread.java:750)`

	srv, cli := testutil.NewMockServer(mockResponse, "/_nodes/hot_threads")
	defer srv.Close()

	shared.SetClient(cli)

	result, err := GetHotThreads(context.Background(), "", 0, "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == "" {
		t.Fatal("expected non-empty result")
	}

	if !strings.Contains(result, "node-1") {
		t.Error("expected result to contain node-1")
	}

	if !strings.Contains(result, "cpu usage") {
		t.Error("expected result to contain 'cpu usage'")
	}
}

func TestGetHotThreadsWithParams(t *testing.T) {
	mockResponse := `Hot threads output`

	// The query parameters will be in the URL
	srv, cli := testutil.NewMockServer(mockResponse, "/_nodes/hot_threads")
	defer srv.Close()

	shared.SetClient(cli)

	result, err := GetHotThreads(context.Background(), "", 5, "1s", "cpu")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Hot threads output" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestGetHotThreadsSpecificNode(t *testing.T) {
	mockResponse := `Hot threads for node1`

	srv, cli := testutil.NewMockServer(mockResponse, "/_nodes/node1/hot_threads")
	defer srv.Close()

	shared.SetClient(cli)

	result, err := GetHotThreads(context.Background(), "node1", 0, "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != "Hot threads for node1" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestGetHotThreadsError(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(500, `Internal server error`, "/_nodes/hot_threads")
	defer srv.Close()

	shared.SetClient(cli)

	_, err := GetHotThreads(context.Background(), "", 0, "", "")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
