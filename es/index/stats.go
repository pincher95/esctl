package index

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pincher95/esctl/shared"
)

// IndexStatsResponse wraps the stats API response
type IndexStatsResponse struct {
	Shards  map[string]any       `json:"_shards"`
	All     IndexStat            `json:"_all"`
	Indices map[string]IndexStat `json:"indices,omitempty"`
}

// IndexStat contains statistics for a single index
type IndexStat struct {
	UUID      string          `json:"uuid,omitempty"`
	Health    string          `json:"health,omitempty"`
	Status    string          `json:"status,omitempty"`
	Primaries IndexStatDetail `json:"primaries"`
	Total     IndexStatDetail `json:"total"`
}

// IndexStatDetail contains detailed statistics
type IndexStatDetail struct {
	Docs     IndexStatDocs     `json:"docs"`
	Store    IndexStatStore    `json:"store"`
	Indexing IndexStatIndexing `json:"indexing"`
	Get      IndexStatGet      `json:"get"`
	Search   IndexStatSearch   `json:"search"`
	Merges   IndexStatMerges   `json:"merges"`
	Refresh  IndexStatRefresh  `json:"refresh"`
	Flush    IndexStatFlush    `json:"flush"`
	Segments IndexStatSegments `json:"segments"`
}

type IndexStatDocs struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

type IndexStatStore struct {
	SizeInBytes int64  `json:"size_in_bytes"`
	Size        string `json:"size,omitempty"`
}

type IndexStatIndexing struct {
	IndexTotal        int64 `json:"index_total"`
	IndexTimeInMillis int64 `json:"index_time_in_millis"`
	IndexCurrent      int64 `json:"index_current"`
}

type IndexStatGet struct {
	Total               int64 `json:"total"`
	TimeInMillis        int64 `json:"time_in_millis"`
	ExistsTotal         int64 `json:"exists_total"`
	ExistsTimeInMillis  int64 `json:"exists_time_in_millis"`
	MissingTotal        int64 `json:"missing_total"`
	MissingTimeInMillis int64 `json:"missing_time_in_millis"`
	Current             int64 `json:"current"`
}

type IndexStatSearch struct {
	QueryTotal        int64 `json:"query_total"`
	QueryTimeInMillis int64 `json:"query_time_in_millis"`
	QueryCurrent      int64 `json:"query_current"`
	FetchTotal        int64 `json:"fetch_total"`
	FetchTimeInMillis int64 `json:"fetch_time_in_millis"`
	FetchCurrent      int64 `json:"fetch_current"`
}

type IndexStatMerges struct {
	Current            int64 `json:"current"`
	CurrentDocs        int64 `json:"current_docs"`
	CurrentSizeInBytes int64 `json:"current_size_in_bytes"`
	Total              int64 `json:"total"`
	TotalTimeInMillis  int64 `json:"total_time_in_millis"`
	TotalDocs          int64 `json:"total_docs"`
	TotalSizeInBytes   int64 `json:"total_size_in_bytes"`
}

type IndexStatRefresh struct {
	Total             int64 `json:"total"`
	TotalTimeInMillis int64 `json:"total_time_in_millis"`
}

type IndexStatFlush struct {
	Total             int64 `json:"total"`
	TotalTimeInMillis int64 `json:"total_time_in_millis"`
}

type IndexStatSegments struct {
	Count                     int64 `json:"count"`
	MemoryInBytes             int64 `json:"memory_in_bytes"`
	TermsMemoryInBytes        int64 `json:"terms_memory_in_bytes"`
	StoredFieldsMemoryInBytes int64 `json:"stored_fields_memory_in_bytes"`
	TermVectorsMemoryInBytes  int64 `json:"term_vectors_memory_in_bytes"`
	NormsMemoryInBytes        int64 `json:"norms_memory_in_bytes"`
	PointsMemoryInBytes       int64 `json:"points_memory_in_bytes"`
	DocValuesMemoryInBytes    int64 `json:"doc_values_memory_in_bytes"`
	IndexWriterMemoryInBytes  int64 `json:"index_writer_memory_in_bytes"`
	VersionMapMemoryInBytes   int64 `json:"version_map_memory_in_bytes"`
	FixedBitSetMemoryInBytes  int64 `json:"fixed_bit_set_memory_in_bytes"`
}

// GetIndexStats retrieves statistics for one or more indices
func (i *index) GetIndexStats(ctx context.Context, indices []string, metric string) (*IndexStatsResponse, error) {
	u := url.URL{}
	if len(indices) > 0 {
		u.Path = fmt.Sprintf("%s/_stats", strings.Join(indices, ","))
	} else {
		u.Path = "_stats"
	}

	if metric != "" {
		u.Path = fmt.Sprintf("%s/%s", u.Path, metric)
	}

	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var out IndexStatsResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get index stats: %s", resp.Status())
	}
	return &out, nil
}

// RecoveryResponse wraps the recovery API response
type RecoveryResponse map[string]IndexRecovery

// IndexRecovery contains recovery information for an index
type IndexRecovery struct {
	Shards []ShardRecovery `json:"shards"`
}

// ShardRecovery contains recovery information for a single shard
type ShardRecovery struct {
	ID                int64               `json:"id"`
	Type              string              `json:"type"`
	Stage             string              `json:"stage"`
	Primary           bool                `json:"primary"`
	StartTime         string              `json:"start_time,omitempty"`
	StartTimeInMillis int64               `json:"start_time_in_millis,omitempty"`
	StopTime          string              `json:"stop_time,omitempty"`
	StopTimeInMillis  int64               `json:"stop_time_in_millis,omitempty"`
	TotalTimeInMillis int64               `json:"total_time_in_millis,omitempty"`
	Source            RecoverySource      `json:"source"`
	Target            RecoveryTarget      `json:"target"`
	Index             RecoveryIndex       `json:"index"`
	Translog          RecoveryTranslog    `json:"translog"`
	VerifyIndex       RecoveryVerifyIndex `json:"verify_index"`
}

type RecoverySource struct {
	ID               string `json:"id,omitempty"`
	Host             string `json:"host,omitempty"`
	TransportAddress string `json:"transport_address,omitempty"`
	IP               string `json:"ip,omitempty"`
	Name             string `json:"name,omitempty"`
}

type RecoveryTarget struct {
	ID               string `json:"id,omitempty"`
	Host             string `json:"host,omitempty"`
	TransportAddress string `json:"transport_address,omitempty"`
	IP               string `json:"ip,omitempty"`
	Name             string `json:"name,omitempty"`
}

type RecoveryIndex struct {
	Size                       RecoverySize  `json:"size"`
	Files                      RecoveryFiles `json:"files"`
	TotalTimeInMillis          int64         `json:"total_time_in_millis"`
	SourceThrottleTimeInMillis int64         `json:"source_throttle_time_in_millis"`
	TargetThrottleTimeInMillis int64         `json:"target_throttle_time_in_millis"`
}

type RecoverySize struct {
	TotalInBytes     int64  `json:"total_in_bytes"`
	ReusedInBytes    int64  `json:"reused_in_bytes"`
	RecoveredInBytes int64  `json:"recovered_in_bytes"`
	Percent          string `json:"percent"`
}

type RecoveryFiles struct {
	Total     int64  `json:"total"`
	Reused    int64  `json:"reused"`
	Recovered int64  `json:"recovered"`
	Percent   string `json:"percent"`
}

type RecoveryTranslog struct {
	Recovered         int64  `json:"recovered"`
	Total             int64  `json:"total"`
	Percent           string `json:"percent"`
	TotalOnStart      int64  `json:"total_on_start"`
	TotalTimeInMillis int64  `json:"total_time_in_millis,omitempty"`
}

type RecoveryVerifyIndex struct {
	CheckIndexTimeInMillis int64 `json:"check_index_time_in_millis"`
	TotalTimeInMillis      int64 `json:"total_time_in_millis"`
}

// GetRecovery retrieves recovery information for one or more indices
func (i *index) GetRecovery(ctx context.Context, indices []string, detailed bool) (RecoveryResponse, error) {
	u := url.URL{}
	if len(indices) > 0 {
		u.Path = fmt.Sprintf("%s/_recovery", strings.Join(indices, ","))
	} else {
		u.Path = "_recovery"
	}

	q := u.Query()
	q.Set("format", "json")
	if detailed {
		q.Set("detailed", "true")
	}
	u.RawQuery = q.Encode()

	var out RecoveryResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get recovery info: %s", resp.Status())
	}
	return out, nil
}

// SegmentsResponse wraps the segments API response
type SegmentsResponse struct {
	Shards  map[string]any           `json:"_shards"`
	Indices map[string]IndexSegments `json:"indices"`
}

// IndexSegments contains segment information for an index
type IndexSegments struct {
	Shards map[string][]ShardSegments `json:"shards"`
}

// ShardSegments contains segment information for a single shard
type ShardSegments struct {
	Routing              ShardRouting           `json:"routing"`
	NumCommittedSegments int                    `json:"num_committed_segments"`
	NumSearchSegments    int                    `json:"num_search_segments"`
	Segments             map[string]SegmentInfo `json:"segments"`
}

type ShardRouting struct {
	State          string `json:"state"`
	Primary        bool   `json:"primary"`
	Node           string `json:"node"`
	RelocatingNode string `json:"relocating_node,omitempty"`
}

type SegmentInfo struct {
	Generation    int64             `json:"generation"`
	NumDocs       int64             `json:"num_docs"`
	DeletedDocs   int64             `json:"deleted_docs"`
	SizeInBytes   int64             `json:"size_in_bytes"`
	MemoryInBytes int64             `json:"memory_in_bytes"`
	Committed     bool              `json:"committed"`
	Search        bool              `json:"search"`
	Version       string            `json:"version"`
	Compound      bool              `json:"compound"`
	Attributes    map[string]string `json:"attributes,omitempty"`
}

// GetSegments retrieves segment information for one or more indices
func (i *index) GetSegments(ctx context.Context, indices []string) (*SegmentsResponse, error) {
	u := url.URL{}
	if len(indices) > 0 {
		u.Path = fmt.Sprintf("%s/_segments", strings.Join(indices, ","))
	} else {
		u.Path = "_segments"
	}

	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var out SegmentsResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get segments: %s", resp.Status())
	}
	return &out, nil
}
