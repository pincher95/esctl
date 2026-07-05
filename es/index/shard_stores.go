package index

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pincher95/esctl/shared"
)

// ShardStoresResponse wraps the _shard_stores API response.
//
// The API reports, for every shard copy that exists on disk, which node holds
// it, its allocation id, whether that copy is currently used as a primary or
// replica (or is unused), and any store-level exception (e.g. corruption).
// It is the primary tool for diagnosing why a shard cannot be allocated.
type ShardStoresResponse struct {
	Indices map[string]ShardStoresIndex `json:"indices"`
}

// ShardStoresIndex holds the per-shard stores for a single index.
type ShardStoresIndex struct {
	Shards map[string]ShardStoresShard `json:"shards"`
}

// ShardStoresShard holds every on-disk copy discovered for one shard.
type ShardStoresShard struct {
	Stores []ShardStore `json:"stores"`
}

// ShardStore describes a single on-disk copy of a shard on a node.
//
// Elasticsearch/OpenSearch encode each store as an object that mixes a dynamic
// "<node-id>" key (whose value is the node info) with the fixed keys
// "allocation_id", "allocation" and "store_exception", so it needs a custom
// unmarshaler to separate the node identity from the fixed fields.
// The response is decoded via a custom UnmarshalJSON (below); the json tags
// here only shape the flattened -o json/yaml output, which is intentionally
// cleaner than the raw API format (no dynamic node-id keys).
type ShardStore struct {
	NodeID           string `json:"node_id"`
	NodeName         string `json:"node_name,omitempty"`
	TransportAddress string `json:"transport_address,omitempty"`
	AllocationID     string `json:"allocation_id"`
	// Allocation is one of "primary", "replica" or "unused".
	Allocation string `json:"allocation"`
	// StoreException is set when the copy cannot be opened (e.g. corruption).
	StoreException *FailuresCause `json:"store_exception,omitempty"`
}

// storeNodeInfo captures the fields we surface from the dynamic node object.
type storeNodeInfo struct {
	Name             string `json:"name"`
	TransportAddress string `json:"transport_address"`
}

// UnmarshalJSON separates the dynamic node-id key from the fixed store fields.
func (s *ShardStore) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		switch key {
		case "allocation_id":
			if err := json.Unmarshal(value, &s.AllocationID); err != nil {
				return err
			}
		case "allocation":
			if err := json.Unmarshal(value, &s.Allocation); err != nil {
				return err
			}
		case "store_exception":
			s.StoreException = &FailuresCause{}
			if err := json.Unmarshal(value, s.StoreException); err != nil {
				return err
			}
		default:
			// Any remaining key is the dynamic node id; its value is the node info.
			var node storeNodeInfo
			if err := json.Unmarshal(value, &node); err != nil {
				return err
			}
			s.NodeID = key
			s.NodeName = node.Name
			s.TransportAddress = node.TransportAddress
		}
	}
	return nil
}

// GetShardStores retrieves shard store information for one or more indices.
// When indices is empty it queries all indices. status filters the returned
// shards by health and is a comma-separated subset of green,yellow,red,all
// (the API default is yellow,red — shards with at least one unassigned copy).
func (i *index) GetShardStores(ctx context.Context, indices []string, status string) (*ShardStoresResponse, error) {
	u := url.URL{}
	if len(indices) > 0 {
		u.Path = fmt.Sprintf("%s/_shard_stores", strings.Join(indices, ","))
	} else {
		u.Path = "_shard_stores"
	}

	q := u.Query()
	q.Set("format", "json")
	if status != "" {
		q.Set("status", status)
	}
	u.RawQuery = q.Encode()

	var out ShardStoresResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get shard stores: %s", resp.Status())
	}
	return &out, nil
}
