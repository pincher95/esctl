package node

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestGetThreadPools(t *testing.T) {
	mockResponse := `{
		"nodes": {
			"node1": {
				"name": "node-1",
				"thread_pool": {
					"search": {
						"threads": 13,
						"queue": 0,
						"active": 0,
						"rejected": 0,
						"largest": 13,
						"completed": 123456
					},
					"write": {
						"threads": 8,
						"queue": 0,
						"active": 0,
						"rejected": 0,
						"largest": 8,
						"completed": 789012
					},
					"get": {
						"threads": 8,
						"queue": 0,
						"active": 0,
						"rejected": 0,
						"largest": 8,
						"completed": 345678
					}
				}
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_nodes/thread_pool")
	defer srv.Close()

	shared.SetClient(cli)

	result, err := GetThreadPools(context.Background(), "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(result.Nodes))
	}

	node := result.Nodes["node1"]
	if node.Name != "node-1" {
		t.Errorf("expected node name 'node-1', got %s", node.Name)
	}

	if len(node.ThreadPools) != 3 {
		t.Errorf("expected 3 thread pools, got %d", len(node.ThreadPools))
	}

	searchPool := node.ThreadPools["search"]
	if searchPool.Threads != 13 {
		t.Errorf("expected search pool threads to be 13, got %d", searchPool.Threads)
	}

	if searchPool.Completed != 123456 {
		t.Errorf("expected search pool completed to be 123456, got %d", searchPool.Completed)
	}
}

func TestGetThreadPoolsSpecificNode(t *testing.T) {
	mockResponse := `{
		"nodes": {
			"node1": {
				"name": "node-1",
				"thread_pool": {
					"search": {
						"threads": 13,
						"queue": 0,
						"active": 0,
						"rejected": 0,
						"largest": 13,
						"completed": 123456
					}
				}
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_nodes/node1/thread_pool")
	defer srv.Close()

	shared.SetClient(cli)

	result, err := GetThreadPools(context.Background(), "node1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(result.Nodes))
	}
}

func TestGetThreadPoolsError(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(500, `{"error":"internal server error"}`, "/_nodes/thread_pool")
	defer srv.Close()

	shared.SetClient(cli)

	_, err := GetThreadPools(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
