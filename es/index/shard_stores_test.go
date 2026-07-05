package index

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

// realistic _shard_stores payload: one healthy primary copy and one unused copy
// carrying a store exception (corruption), using the dynamic "<node-id>" key.
const shardStoresJSON = `{
  "indices": {
    "my-index": {
      "shards": {
        "0": {
          "stores": [
            {
              "node-abc": {
                "name": "node_t0",
                "transport_address": "127.0.0.1:9300",
                "attributes": {}
              },
              "allocation_id": "alloc-123",
              "allocation": "primary"
            },
            {
              "node-def": {
                "name": "node_t1",
                "transport_address": "127.0.0.1:9301"
              },
              "allocation_id": "alloc-456",
              "allocation": "unused",
              "store_exception": {
                "type": "corrupted_index_exception",
                "reason": "failed engine on shard"
              }
            }
          ]
        }
      }
    }
  }
}`

func TestGetShardStores(t *testing.T) {
	srv, cli := testutil.NewMockServer(shardStoresJSON, "/my-index/_shard_stores")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetShardStores(context.Background(), []string{"my-index"}, "red,yellow")
	if err != nil {
		t.Fatalf("GetShardStores() error = %v", err)
	}

	index, ok := resp.Indices["my-index"]
	if !ok {
		t.Fatalf("expected index 'my-index' in response, got %v", resp.Indices)
	}
	shard, ok := index.Shards["0"]
	if !ok {
		t.Fatalf("expected shard '0' in response, got %v", index.Shards)
	}
	if len(shard.Stores) != 2 {
		t.Fatalf("expected 2 stores, got %d", len(shard.Stores))
	}

	primary := shard.Stores[0]
	if primary.NodeID != "node-abc" {
		t.Errorf("expected NodeID 'node-abc', got %q", primary.NodeID)
	}
	if primary.NodeName != "node_t0" {
		t.Errorf("expected NodeName 'node_t0', got %q", primary.NodeName)
	}
	if primary.TransportAddress != "127.0.0.1:9300" {
		t.Errorf("expected TransportAddress '127.0.0.1:9300', got %q", primary.TransportAddress)
	}
	if primary.AllocationID != "alloc-123" {
		t.Errorf("expected AllocationID 'alloc-123', got %q", primary.AllocationID)
	}
	if primary.Allocation != "primary" {
		t.Errorf("expected Allocation 'primary', got %q", primary.Allocation)
	}
	if primary.StoreException != nil {
		t.Errorf("expected no store exception on primary, got %+v", primary.StoreException)
	}

	corrupt := shard.Stores[1]
	if corrupt.NodeID != "node-def" {
		t.Errorf("expected NodeID 'node-def', got %q", corrupt.NodeID)
	}
	if corrupt.Allocation != "unused" {
		t.Errorf("expected Allocation 'unused', got %q", corrupt.Allocation)
	}
	if corrupt.StoreException == nil {
		t.Fatalf("expected a store exception on the corrupt copy, got nil")
	}
	if corrupt.StoreException.Type != "corrupted_index_exception" {
		t.Errorf("expected exception type 'corrupted_index_exception', got %q", corrupt.StoreException.Type)
	}
}

func TestGetShardStoresAllIndices(t *testing.T) {
	srv, cli := testutil.NewMockServer(`{"indices":{}}`, "/_shard_stores")
	defer srv.Close()
	shared.SetClient(cli)

	idx := NewIndex()
	resp, err := idx.GetShardStores(context.Background(), nil, "")
	if err != nil {
		t.Fatalf("GetShardStores() error = %v", err)
	}
	if len(resp.Indices) != 0 {
		t.Errorf("expected 0 indices, got %d", len(resp.Indices))
	}
}
