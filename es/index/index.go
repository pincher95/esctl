package index

import "context"

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

// ResponseShards is a sub type of api repsonses containing information about shards
type ResponseShards struct {
	Total      int                     `json:"total"`
	Successful int                     `json:"successful"`
	Failed     int                     `json:"failed"`
	Failures   []ResponseShardsFailure `json:"failures"`
	Skipped    int                     `json:"skipped"`
}

// ResponseShardsFailure is a sub type of ReponseShards containing information about a failed shard
type ResponseShardsFailure struct {
	Shard  int    `json:"shard"`
	Index  any    `json:"index"`
	Node   string `json:"node"`
	Reason struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"reason"`
}

// FailuresCause contains information about failure cause
type FailuresCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	NodeID string `json:"node_id"`
	Cause  *struct {
		Type   string  `json:"type"`
		Reason *string `json:"reason"`
	} `json:"caused_by,omitempty"`
}

// FailuresShard contains information about shard failures
type FailuresShard struct {
	Shard  int           `json:"shard"`
	Index  string        `json:"index"`
	Status string        `json:"status"`
	Reason FailuresCause `json:"reason"`
}
