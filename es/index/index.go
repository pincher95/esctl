package index

import (
	"context"

	"github.com/pincher95/esctl/es/common"
)

type IndexSettingsResponse map[string]any

// Index exposes typed helpers around index-level Elasticsearch APIs.
// All parameters are plain values; pass empty strings when you don’t want to set them.
type Index interface {
	UpdateIndexSettings(ctx context.Context, index string, body map[string]any, flatSettings bool) (*IndexSettingsResponse, error)
	GetAliases(ctx context.Context, index string) (*IndexAliasResponse, error)
	CacheClear(ctx context.Context, index string) (*IndexCacheClearResponse, error)
}

type index struct{}

func NewIndex() Index {
	return &index{}
}

// ResponseShards is shared via es/common.
type ResponseShards = common.ResponseShards

// ResponseShardsFailure is shared via es/common.
type ResponseShardsFailure = common.ResponseShardsFailure

// FailuresCause is shared via es/common to avoid duplication across packages.
type FailuresCause = common.FailuresCause

// FailuresShard is shared via es/common.
type FailuresShard = common.FailuresShard
