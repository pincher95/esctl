package cluster

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type AllocationExplain struct {
	Index          string                       `json:"index"`
	Shard          int                          `json:"shard"`
	Primary        bool                         `json:"primary"`
	CurrentState   string                       `json:"current_state"`
	CurrentNode    ClusterAllocationCurrentNode `json:"current_node"`
	UnassignedInfo struct {
		Reason               string `json:"reason"`
		At                   string `json:"at"`
		LastAllocationStatus string `json:"last_allocation_status"`
	} `json:"unassigned_info"`
	CanAllocate                  string                             `json:"can_allocate"`
	CanRemainOnCurrentNode       string                             `json:"can_remain_on_current_node"`
	CanRebalanceCluster          string                             `json:"can_rebalance_cluster"`
	CanRebalanceToOtherNode      string                             `json:"can_rebalance_to_other_node"`
	RebalanceExplanation         string                             `json:"rebalance_explanation"`
	AllocateExplanation          string                             `json:"allocate_explanation"`
	NodeAllocationDecisions      []ClusterAllocationNodeDecisions   `json:"node_allocation_decisions"`
	CanRebalanceClusterDecisions []ClusterAllocationExplainDeciders `json:"can_rebalance_cluster_decisions,omitempty"`
}

// ClusterAllocationCurrentNode is a sub type of ClusterAllocationExplainResp containing information of the node the shard is on
type ClusterAllocationCurrentNode struct {
	NodeID           string `json:"id"`
	NodeName         string `json:"name"`
	TransportAddress string `json:"transport_address"`
	NodeAttributes   struct {
		ShardIndexingPressureEnabled string `json:"shard_indexing_pressure_enabled"`
	} `json:"attributes"`
	WeightRanking int `json:"weight_ranking"`
}

// ClusterAllocationNodeDecisions is a sub type of ClusterAllocationExplainResp containing information of a node allocation decission
type ClusterAllocationNodeDecisions struct {
	NodeID           string `json:"node_id"`
	NodeName         string `json:"node_name"`
	TransportAddress string `json:"transport_address"`
	NodeAttributes   struct {
		ShardIndexingPressureEnabled string `json:"shard_indexing_pressure_enabled"`
	} `json:"node_attributes"`
	NodeDecision  string                             `json:"node_decision"`
	WeightRanking int                                `json:"weight_ranking"`
	Deciders      []ClusterAllocationExplainDeciders `json:"deciders"`
}

// ClusterAllocationExplainDeciders is a sub type of ClusterAllocationExplainResp and
// ClusterAllocationNodeDecisions containing inforamtion about Deciders decissions
type ClusterAllocationExplainDeciders struct {
	Decider     string `json:"decider"`
	Decision    string `json:"decision"`
	Explanation string `json:"explanation"`
}

func ClusterAllocationExplain(ctx context.Context, includeDiskInfo, includeYesDecisions bool) (*AllocationExplain, error) {
	u := url.URL{Path: "_cluster/allocation/explain"}
	q := u.Query()
	q.Set("format", "json")
	if includeDiskInfo {
		q.Set("include_disk_info", "")
	}
	if includeYesDecisions {
		q.Set("include_yes_decisions", "")
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out AllocationExplain
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%s", resp.String())
	}
	return &out, nil
}
