package cat

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type CatAllocationResponse struct {
	Shards int `json:"shards,string"`
	// API returns null for certain fields when node is UNASSIGNED
	DiskIndices *string `json:"disk.indices"`
	DiskUsed    *string `json:"disk.used"`
	DiskAvail   *string `json:"disk.avail"`
	DiskTotal   *string `json:"disk.total"`
	DiskPercent *int    `json:"disk.percent,string"`
	Host        *string `json:"host"`
	IP          *string `json:"ip"`
	Node        string  `json:"node"`
}

func (c *cat) CatAllocation(ctx context.Context, endpoint, nodeID, bytes string) ([]CatAllocationResponse, error) {
	// Allow caller to override full endpoint when necessary
	if endpoint == "" {
		path := "_cat/allocation"
		if nodeID != "" {
			path = fmt.Sprintf("_cat/allocation/%s", nodeID)
		}

		u := url.URL{Path: path}
		q := u.Query()
		q.Set("format", "json")
		q.Set("h", "shards,disk.indices,disk.used,disk.avail,disk.total,host,ip,node,disk.percent")
		if bytes != "" {
			q.Set("bytes", bytes)
		}
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	allocations := make([]CatAllocationResponse, 0)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&allocations).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get nodes allocations: %s", resp.Status())
	}

	return allocations, nil
}
