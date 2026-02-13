package index

import (
	"context"

	"github.com/pincher95/esctl/es/common"
)

type IndexSettingsResponse map[string]any

// Index exposes typed helpers around index-level Elasticsearch APIs.
// All parameters are plain values; pass empty strings when you don't want to set them.
type Index interface {
	UpdateIndexSettings(ctx context.Context, index string, body map[string]any, flatSettings bool) (*IndexSettingsResponse, error)
	GetAliases(ctx context.Context, index string) (*AliasListResponse, error)
	CacheClear(ctx context.Context, index string) (*IndexCacheClearResponse, error)
	Refresh(ctx context.Context, index string) (*IndexOpResponse, error)
	Flush(ctx context.Context, index string) (*IndexOpResponse, error)
	Forcemerge(ctx context.Context, index string, maxNumSegments int, onlyExpungeDeletes bool, flush bool) (*IndexOpResponse, error)
	// Lifecycle operations
	Open(ctx context.Context, indices []string) (*IndexOperationResponse, error)
	Close(ctx context.Context, indices []string) (*IndexOperationResponse, error)
	Clone(ctx context.Context, source, target string, settings map[string]any) (*IndexOperationResponse, error)
	Split(ctx context.Context, source, target string, shards int, settings map[string]any) (*IndexOperationResponse, error)
	Shrink(ctx context.Context, source, target string, shards int, settings map[string]any) (*IndexOperationResponse, error)
	// Stats and diagnostics
	GetIndexStats(ctx context.Context, indices []string, metric string) (*IndexStatsResponse, error)
	GetRecovery(ctx context.Context, indices []string, detailed bool) (RecoveryResponse, error)
	GetSegments(ctx context.Context, indices []string) (*SegmentsResponse, error)
}

type index struct{}

// NewIndex constructs an Index implementation backed by the default HTTP client.
func NewIndex() Index {
	return &index{}
}

// GetAliases implements the Index interface
func (i *index) GetAliases(ctx context.Context, indexName string) (*AliasListResponse, error) {
	indices := []string{}
	if indexName != "" {
		indices = []string{indexName}
	}
	result, err := GetAlias(ctx, indices, []string{})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ResponseShards is shared via es/common.
type ResponseShards = common.ResponseShards

// ResponseShardsFailure is shared via es/common.
type ResponseShardsFailure = common.ResponseShardsFailure

// FailuresCause is shared via es/common to avoid duplication across packages.
type FailuresCause = common.FailuresCause

// FailuresShard is shared via es/common.
type FailuresShard = common.FailuresShard
