package node

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// GetHotThreads retrieves hot threads information for nodes
// threadType can be: "cpu", "wait", or "block"
// interval is a time value like "500ms" or "1s"
// threads is the number of hot threads to return (default 3)
func GetHotThreads(ctx context.Context, nodeID string, threads int, interval string, threadType string) (string, error) {
	endpoint := "/_nodes/hot_threads"
	if nodeID != "" {
		endpoint = fmt.Sprintf("/_nodes/%s/hot_threads", nodeID)
	}

	// Build query parameters
	params := url.Values{}
	if threads > 0 {
		params.Add("threads", fmt.Sprintf("%d", threads))
	}
	if interval != "" {
		params.Add("interval", interval)
	}
	if threadType != "" {
		params.Add("type", threadType)
	}

	if len(params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		Get(endpoint)

	if err != nil {
		return "", fmt.Errorf("failed to get hot threads: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("error getting hot threads: %s", resp.Status())
	}

	// Hot threads returns plain text, not JSON
	return resp.String(), nil
}
