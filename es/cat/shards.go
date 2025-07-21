package cat

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type CatShardResponse struct {
	ShardIdentity       // embedded: index, shard, prirep, state
	ShardStorageStats   // embedded: docs, store
	ShardNodeInfo       // embedded: ip, id, node
	ShardUnassignedInfo // unassigned.* fields

	SyncID                         *string `json:"sync_id"`
	CompletionSize                 *string `json:"completion.size"`
	FielddataMemorySize            *string `json:"fielddata.memory_size"`
	FielddataEvictions             *int    `json:"fielddata.evictions,string"`
	QueryCacheMemorySize           *string `json:"query_cache.memory_size"`
	QueryCacheEvictions            *int    `json:"query_cache.evictions,string"`
	FlushTotal                     *int    `json:"flush.total,string"`
	FlushTotalTime                 *string `json:"flush.total_time"`
	GetCurrent                     *int    `json:"get.current,string"`
	GetTime                        *string `json:"get.time"`
	GetTotal                       *int    `json:"get.total,string"`
	GetExistsTime                  *string `json:"get.exists_time"`
	GetExistsTotal                 *int    `json:"get.exists_total,string"`
	GetMissingTime                 *string `json:"get.missing_time"`
	GetMissingTotal                *int    `json:"get.missing_total,string"`
	IndexingDeleteCurrent          *int    `json:"indexing.delete_current,string"`
	IndexingDeleteTime             *string `json:"indexing.delete_time"`
	IndexingDeleteTotal            *string `json:"indexing.delete_total"`
	IndexingIndexCurrent           *int    `json:"indexing.index_current,string"`
	IndexingIndexTime              *string `json:"indexing.index_time"`
	IndexingIndexTotal             *int    `json:"indexing.index_total,string"`
	IndexingIndexFailed            *int    `json:"indexing.index_failed,string"`
	MergesCurrent                  *int    `json:"merges.current,string"`
	MergesCurrentDocs              *int    `json:"merges.current_docs,string"`
	MergesCurrentSize              *string `json:"merges.current_size"`
	MergesTotal                    *int    `json:"merges.total,string"`
	MergesTotalDocs                *int    `json:"merges.total_docs,string"`
	MergesTotalSize                *string `json:"merges.total_size"`
	MergesTotalTime                *string `json:"merges.total_time"`
	RefreshTotal                   *int    `json:"refresh.total,string"`
	RefreshTime                    *string `json:"refresh.time"`
	RefreshExternalTotal           *int    `json:"refresh.external_total,string"`
	RefreshExternalTime            *string `json:"refresh.external_time"`
	RefreshListeners               *int    `json:"refresh.listeners,string"`
	SearchFetchCurrent             *int    `json:"search.fetch_current,string"`
	SearchFetchTime                *string `json:"search.fetch_time"`
	SearchFetchTotal               *int    `json:"search.fetch_total,string"`
	SearchOpenContexts             *int    `json:"search.open_contexts,string"`
	SearchQueryCurrent             *int    `json:"search.query_current,string"`
	SearchQueryTime                *string `json:"search.query_time"`
	SearchQueryTotal               *int    `json:"search.query_total,string"`
	SearchConcurrentQueryCurrent   *int    `json:"search.concurrent_query_current,string"`
	SearchConcurrentQueryTime      *string `json:"search.concurrent_query_time"`
	SearchConcurrentQueryTotal     *int    `json:"search.concurrent_query_total,string"`
	SearchConcurrentAvgSliceCount  *string `json:"search.concurrent_avg_slice_count"`
	SearchScrollCurrent            *int    `json:"search.scroll_current,string"`
	SearchScrollTime               *string `json:"search.scroll_time"`
	SearchScrollTotal              *int    `json:"search.scroll_total,string"`
	SearchPointInTimeCurrent       *int    `json:"search.point_in_time_current,string"`
	SearchPointInTimeTime          *string `json:"search.point_in_time_time"`
	SearchPointInTimeTotal         *int    `json:"search.point_in_time_total,string"`
	SearchIdleReactivateCountTotal *int    `json:"search.search_idle_reactivate_count_total,string"`
	SegmentsCount                  *int    `json:"segments.count,string"`
	SegmentsMemory                 *string `json:"segments.memory"`
	SegmentsIndexWriterMemory      *string `json:"segments.index_writer_memory"`
	SegmentsVersionMapMemory       *string `json:"segments.version_map_memory"`
	SegmentsFixedBitsetMemory      *string `json:"segments.fixed_bitset_memory"`
	SeqNoMax                       *int    `json:"seq_no.max,string"`
	SeqNoLocalCheckpoint           *int    `json:"seq_no.local_checkpoint,string"`
	SeqNoGlobalCheckpoint          *int    `json:"seq_no.global_checkpoint,string"`
	WarmerCurrent                  *int    `json:"warmer.current,string"`
	WarmerTotal                    *int    `json:"warmer.total,string"`
	WarmerTotalTime                *string `json:"warmer.total_time"`
	PathData                       *string `json:"path.data"`
	PathState                      *string `json:"path.state"`
	DocsDeleted                    *int    `json:"docs.deleted,string"`
}

func (c *cat) CatShards(ctx context.Context, endpoint, index, bytes, timeUnit string) ([]CatShardResponse, error) {
	if endpoint == "" {
		path := "_cat/shards"
		if index != "" {
			path = fmt.Sprintf("_cat/shards/%s", index)
		}

		u := url.URL{Path: path}
		q := u.Query()
		q.Set("format", "json")
		q.Set("h", "index,shard,prirep,state,docs,store,ip,id,node,unassigned.reason,unassigned.at,segments.count")
		if bytes != "" {
			q.Set("bytes", bytes)
		}
		if timeUnit != "" {
			q.Set("time", timeUnit)
		}
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	shards := make([]CatShardResponse, 0)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&shards).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get shards: %s", resp.Status())
	}

	return shards, nil
}
