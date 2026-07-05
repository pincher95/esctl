package node

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// NodeStats surfaces the health-critical subset of the _nodes/stats response.
type NodeStats struct {
	Name string `json:"name"`
	JVM  struct {
		Mem struct {
			HeapUsedPercent int   `json:"heap_used_percent"`
			HeapUsedInBytes int64 `json:"heap_used_in_bytes"`
			HeapMaxInBytes  int64 `json:"heap_max_in_bytes"`
		} `json:"mem"`
		GC struct {
			Collectors struct {
				Young NodeGCCollector `json:"young"`
				Old   NodeGCCollector `json:"old"`
			} `json:"collectors"`
		} `json:"gc"`
	} `json:"jvm"`
	ThreadPool map[string]struct {
		Queue    int   `json:"queue"`
		Active   int   `json:"active"`
		Rejected int64 `json:"rejected"`
	} `json:"thread_pool"`
	Breakers map[string]struct {
		Tripped int64 `json:"tripped"`
	} `json:"breakers"`
	FS struct {
		Total struct {
			TotalInBytes     int64 `json:"total_in_bytes"`
			AvailableInBytes int64 `json:"available_in_bytes"`
		} `json:"total"`
	} `json:"fs"`
}

// NodeGCCollector holds garbage-collection counters for a single collector.
type NodeGCCollector struct {
	CollectionCount        int64 `json:"collection_count"`
	CollectionTimeInMillis int64 `json:"collection_time_in_millis"`
}

// NodeStatsResponse wraps the _nodes/stats response keyed by node id.
type NodeStatsResponse struct {
	Nodes map[string]NodeStats `json:"nodes"`
}

// GetNodeStats retrieves node statistics (jvm, thread pools, circuit breakers, fs)
// for one node or all nodes. It requests only the health-relevant metric groups.
func GetNodeStats(ctx context.Context, nodeID string) (*NodeStatsResponse, error) {
	path := "_nodes/stats/jvm,thread_pool,breaker,fs"
	if nodeID != "" {
		path = fmt.Sprintf("_nodes/%s/stats/jvm,thread_pool,breaker,fs", nodeID)
	}
	u := url.URL{Path: path}
	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var out NodeStatsResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get node stats: %s", resp.Status())
	}
	return &out, nil
}
