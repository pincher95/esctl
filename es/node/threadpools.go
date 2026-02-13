package node

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

// ThreadPoolStats represents statistics for a thread pool
type ThreadPoolStats struct {
	Threads   int `json:"threads"`
	Queue     int `json:"queue"`
	Active    int `json:"active"`
	Rejected  int `json:"rejected"`
	Largest   int `json:"largest"`
	Completed int `json:"completed"`
}

// NodeThreadPools represents thread pool stats for a single node
type NodeThreadPools struct {
	Name        string                     `json:"name"`
	ThreadPools map[string]ThreadPoolStats `json:"thread_pools"`
}

// ThreadPoolsResponse represents the response from the thread pool API
type ThreadPoolsResponse struct {
	Nodes map[string]struct {
		Name        string                     `json:"name"`
		ThreadPools map[string]ThreadPoolStats `json:"thread_pool"`
	} `json:"nodes"`
}

// GetThreadPools retrieves thread pool statistics for nodes
func GetThreadPools(ctx context.Context, nodeID string) (*ThreadPoolsResponse, error) {
	endpoint := "/_nodes/thread_pool"
	if nodeID != "" {
		endpoint = fmt.Sprintf("/_nodes/%s/thread_pool", nodeID)
	}

	var result ThreadPoolsResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get thread pools: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error getting thread pools: %s", resp.Status())
	}

	return &result, nil
}
