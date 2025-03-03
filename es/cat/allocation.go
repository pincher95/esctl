package cat

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type Allocation struct {
	Shards int `json:"shards,string"`
	// Pointer of string as the api can returns null for those fileds with Node set to "UNASSIGNED"
	DiskIndices *string `json:"disk.indices"`
	DiskUsed    *string `json:"disk.used"`
	DiskAvail   *string `json:"disk.avail"`
	DiskTotal   *string `json:"disk.total"`
	DiskPercent *int    `json:"disk.percent,string"`
	Host        *string `json:"host"`
	IP          *string `json:"ip"`
	Node        string  `json:"node"`
}

func CatAllocation(endpoint, nodeID, bytes *string, debug bool) ([]Allocation, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cat/allocation?format=json&h=shards,disk.indices,disk.used,disk.avail,disk.total,host,ip,node,disk.percent"

		if *nodeID != "" {
			*endpoint = fmt.Sprintf("_cat/allocation/%s?format=json&h=shards,disk.indices,disk.used,disk.avail,disk.total,host,ip,node,disk.percent", *nodeID)
		}
	}

	if bytes != nil {
		*endpoint += fmt.Sprintf("&bytes=%s", *bytes)
	}

	allocations := make([]Allocation, 0)

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&allocations).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get nodes allocations: %s", resp.Status())
	}

	return allocations, nil
}
