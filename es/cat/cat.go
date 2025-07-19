package cat

import "context"

// Cat exposes typed wrappers around the Elasticsearch _cat API.
// All parameters are plain values; pass empty strings when you don’t want to set them.
type Cat interface {
	CatAllocation(ctx context.Context, endpoint, nodeID, bytes string) ([]CatAllocationResponse, error)
	CatFielddata(ctx context.Context, endpoint, fields, bytes string) ([]CatFielddataResponse, error)
	CatIndices(ctx context.Context, endpoint, index, bytes string) ([]CatIndiceResponse, error)
	CatNodes(ctx context.Context, endpoint, nodeName, bytes, timeUnit string) ([]CatNodesResponse, error)
	CatShards(ctx context.Context, endpoint, index, bytes, timeUnit string) ([]CatShardResponse, error)
	CatSnapshots(ctx context.Context, endpoint, repository, snapshotName string) ([]CatSnapshotResponse, error)
	CatPlugins(ctx context.Context, endpoint string) ([]CatPluginResponse, error)
}

type cat struct{}

func NewCat() Cat {
	return &cat{}
}
