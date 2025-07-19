package cluster

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type Health struct {
	ClusterName                 string                 `json:"cluster_name"`
	Status                      string                 `json:"status"`
	TimedOut                    bool                   `json:"timed_out"`
	NumberOfNodes               int                    `json:"number_of_nodes"`
	NumberOfDataNodes           int                    `json:"number_of_data_nodes"`
	DiscoveredMaster            bool                   `json:"discovered_master"`
	DiscoveredClusterManager    bool                   `json:"discovered_cluster_manager"`
	ActivePrimaryShards         int                    `json:"active_primary_shards"`
	ActiveShards                int                    `json:"active_shards"`
	RelocatingShards            int                    `json:"relocating_shards"`
	InitializingShards          int                    `json:"initializing_shards"`
	UnassignedShards            int                    `json:"unassigned_shards"`
	DelayedUnassignedShards     int                    `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int                    `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int                    `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int                    `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64                `json:"active_shards_percent_as_number"`
	Indices                     map[string]IndexHealth `json:"indices,omitempty"`
}

// IndexHealth matches each index's health object
type IndexHealth struct {
	ActivePrimaryShards     int    `json:"active_primary_shards"`
	ActiveShards            int    `json:"active_shards"`
	InitializingShards      int    `json:"initializing_shards"`
	NumberOfReplicas        int    `json:"number_of_replicas"`
	NumberOfShards          int    `json:"number_of_shards"`
	RelocatingShards        int    `json:"relocating_shards"`
	UnassignedShards        int    `json:"unassigned_shards"`
	UnassingedPrimaryShards int    `json:"unassinged_primary_shards"`
	Status                  string `json:"status"`

	// Shards is also a map: shard-id → shard-health
	Shards map[string]ShardHealth `json:"shards,omitempty"`
}

// ShardHealth matches each shard’s health object
type ShardHealth struct {
	ActiveShards            int    `json:"active_shards"`
	InitializingShards      int    `json:"initializing_shards"`
	PrimaryActive           bool   `json:"primary_active"`
	RelocatingShards        int    `json:"relocating_shards"`
	Status                  string `json:"status"`
	UnassignedShards        int    `json:"unassigned_shards"`
	UnassingedPrimaryShards int    `json:"unassinged_primary_shards"`
}

func ClusterHealth(ctx context.Context, level, expandWildcard, index string) (*Health, error) {
	u := url.URL{}
	if index != "" {
		u.Path = fmt.Sprintf("_cluster/health/%s", index)
	} else {
		u.Path = "_cluster/health"
	}
	q := u.Query()
	q.Set("format", "json")
	if level != "" {
		q.Set("level", level)
	}
	if expandWildcard != "" {
		q.Set("expand_wildcards", expandWildcard)
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var result Health
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get cluster health: %s", resp.Status())
	}
	return &result, nil
}
