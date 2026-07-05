package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type Reroute struct {
	Acknowledged bool            `json:"acknowledged"`
	State        json.RawMessage `json:"state"`
	Explanations json.RawMessage `json:"explanations"`
}

// ClusterRerouteState is a sub type of ClusterRerouteResp containing information about the cluster and cluster routing
type ClusterRerouteState struct {
	ClusterUUID        string                       `json:"cluster_uuid"`
	Version            int                          `json:"version"`
	StateUUID          string                       `json:"state_uuid"`
	MasterNode         string                       `json:"master_node"`
	ClusterManagerNode string                       `json:"cluster_manager_node"`
	Blocks             json.RawMessage              `json:"blocks"`
	Nodes              map[string]ClusterStateNodes `json:"nodes"`
	RoutingTable       struct {
		Indices map[string]struct {
			Shards map[string][]ClusterStateRoutingIndex `json:"shards"`
		} `json:"indices"`
	} `json:"routing_table"`
	RoutingNodes      ClusterStateRoutingNodes `json:"routing_nodes"`
	RepositoryCleanup struct {
		RepositoryCleanup []json.RawMessage `json:"repository_cleanup"`
	} `json:"repository_cleanup"`
	SnapshotDeletions struct {
		SnapshotDeletions []json.RawMessage `json:"snapshot_deletions"`
	} `json:"snapshot_deletions"`
	Snapshots struct {
		Snapshots []json.RawMessage `json:"snapshots"`
	} `json:"snapshots"`
	Restore struct {
		Snapshots []json.RawMessage `json:"snapshots"`
	} `json:"restore"`
}

// ClusterStateNodes is a sub type of ClusterStateResp
type ClusterStateNodes struct {
	Name             string            `json:"name"`
	EphemeralID      string            `json:"ephemeral_id"`
	TransportAddress string            `json:"transport_address"`
	Attributes       map[string]string `json:"attributes"`
}

// ClusterStateRoutingIndex is a sub type of ClusterStateResp and ClusterStateRoutingNodes containing information about shard routing
type ClusterStateRoutingIndex struct {
	State                    string  `json:"state"`
	Primary                  bool    `json:"primary"`
	Node                     *string `json:"node"`
	RelocatingNode           *string `json:"relocating_node"`
	Shard                    int     `json:"shard"`
	Index                    string  `json:"index"`
	ExpectedShardSizeInBytes int     `json:"expected_shard_size_in_bytes"`
	AllocationID             *struct {
		ID string `json:"id"`
	} `json:"allocation_id,omitempty"`
	RecoverySource *struct {
		Type string `json:"type"`
	} `json:"recovery_source,omitempty"`
	UnassignedInfo *struct {
		Reason           string `json:"reason"`
		At               string `json:"at"`
		Delayed          bool   `json:"delayed"`
		AllocationStatus string `json:"allocation_status"`
		Details          string `json:"details"`
	} `json:"unassigned_info,omitempty"`
}

// ClusterStateRoutingNodes is a sub type of ClusterStateResp containing information about shard assigned to nodes
type ClusterStateRoutingNodes struct {
	Unassigned []ClusterStateRoutingIndex            `json:"unassigned"`
	Nodes      map[string][]ClusterStateRoutingIndex `json:"nodes"`
}

// ClusterRerouteCommands executes explicit reroute commands (allocate_stale_primary,
// allocate_empty_primary, move, cancel) against the cluster. Each command is a
// single-key map as accepted by the _cluster/reroute API. When dryRun is true the
// cluster state is only simulated and returned, not applied.
func ClusterRerouteCommands(ctx context.Context, commands []map[string]any, dryRun, explain bool) (*Reroute, error) {
	u := url.URL{Path: "_cluster/reroute"}
	q := u.Query()
	q.Set("format", "json")
	q.Set("dry_run", fmt.Sprintf("%t", dryRun))
	q.Set("explain", fmt.Sprintf("%t", explain))
	u.RawQuery = q.Encode()

	body := map[string]any{"commands": commands}

	var out Reroute
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&out).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to post cluster reroute: %s - %s", resp.Status(), string(resp.Body()))
	}
	return &out, nil
}

func ClusterReroute(ctx context.Context, metric string, dryRun, explain, retryFailed bool) (*Reroute, error) {
	u := url.URL{Path: "_cluster/reroute"}
	q := u.Query()
	q.Set("format", "json")
	if metric != "" {
		q.Set("metric", metric)
	}
	q.Set("dry_run", fmt.Sprintf("%t", dryRun))
	q.Set("explain", fmt.Sprintf("%t", explain))
	q.Set("retry_failed", fmt.Sprintf("%t", retryFailed))
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out Reroute
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to post cluster reroute: %s", resp.Status())
	}
	return &out, nil
}
