package cluster

import (
	"encoding/json"
	"fmt"

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

func ClusterReroute(endpoint, flagMertic *string, dryRun, explain, retryFailed bool) (*Reroute, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cluster/reroute?format=json"
	}

	if dryRun {
		*endpoint += fmt.Sprintf("&dry_run=%t", dryRun)
	}

	if explain {
		*endpoint += fmt.Sprintf("&explain=%t", explain)
	}

	if retryFailed {
		*endpoint += fmt.Sprintf("&retry_failed=%t", retryFailed)
	}

	if *flagMertic != "" {
		*endpoint += fmt.Sprintf("&metric=%s", *flagMertic)
	}

	var reroute Reroute

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&reroute).Post(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to post cluster reroute: %s", resp.Status())
	}

	return &reroute, nil
}
