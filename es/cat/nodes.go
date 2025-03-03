package cat

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type Node struct {
	ID                              string  `json:"id"`
	PID                             *string `json:"pid"`
	IP                              string  `json:"ip"`
	Port                            int     `json:"port,string"`
	HTTPAddress                     string  `json:"http_address"`
	Version                         string  `json:"version"`
	Type                            *string `json:"type"`
	Build                           *string `json:"build"`
	JDK                             *string `json:"jdk"`
	DiskTotal                       *string `json:"disk.total"`
	DiskUsed                        *string `json:"disk.used"`
	DiskAvail                       *string `json:"disk.avail"`
	DiskUsedPercent                 *string `json:"disk.used_percent"`
	HeapCurrent                     *string `json:"heap.current"`
	HeapPercent                     *int    `json:"heap.percent,string"`
	HeapMax                         *string `json:"heap.max"`
	RAMCurrent                      *string `json:"ram.current"`
	RAMPercent                      *int    `json:"ram.percent,string"`
	RAMMax                          *string `json:"ram.max"`
	FileDescCurrent                 *int    `json:"file_desc.current,string"`
	FileDescPercent                 *int    `json:"file_desc.percent,string"`
	FileDescMax                     *int    `json:"file_desc.max,string"`
	CPU                             *int    `json:"cpu,string"`
	Load1M                          *string `json:"load_1m"`
	Load5M                          *string `json:"load_5m"`
	Load15M                         *string `json:"load_15m"`
	Uptime                          *string `json:"uptime"`
	Role                            string  `json:"node.role"`
	Roles                           string  `json:"node.roles"`
	Master                          string  `json:"master"`
	ClusterManager                  string  `json:"cluster_manager"`
	Name                            string  `json:"name"`
	CompletionSize                  *string `json:"completion.size"`
	FieldDataMemorySize             *string `json:"fielddata.memory_size"`
	FileldDataEvictions             *int    `json:"fielddata.evictions,string"`
	QueryCacheMemorySize            *string `json:"query_cache.memory_size"`
	QueryCacheEvictions             *int    `json:"query_cache.evictions,string"`
	QueryCacheHitCount              *int    `json:"query_cache.hit_count,string"`
	QueryCacheMissCount             *int    `json:"query_cache.miss_count,string"`
	RequestCacheMemorySize          *string `json:"request_cache.memory_size"`
	RequestCacheEvictions           *int    `json:"request_cache.evictions,string"`
	RequestCacheHitCount            *int    `json:"request_cache.hit_count,string"`
	RequestCacheMissCount           *int    `json:"request_cache.miss_count,string"`
	FlushTotal                      *int    `json:"flush.total,string"`
	FlushTotalTime                  *string `json:"flush.total_time"`
	GetCurrent                      *int    `json:"get.current,string"`
	GetTime                         *string `json:"get.time"`
	GetTotal                        *int    `json:"get.total,string"`
	GetExistsTime                   *string `json:"get.exists_time"`
	GetExistsTotal                  *int    `json:"get.exists_total,string"`
	GetMissingTime                  *string `json:"get.missing_time"`
	GetMissingTotal                 *int    `json:"get.missing_total,string"`
	IndexingDeleteCurrent           *int    `json:"indexing.delete_current,string"`
	IndexingDeleteTime              *string `json:"indexing.delete_time"`
	IndexingDeleteTotal             *int    `json:"indexing.delete_total,string"`
	IndexingIndexCurrent            *int    `json:"indexing.index_current,string"`
	IndexingIndexTime               *string `json:"indexing.index_time"`
	IndexingIndexTotal              *int    `json:"indexing.index_total,string"`
	IndexingIndexFailed             *int    `json:"indexing.index_failed,string"`
	MergesCurrent                   *int    `json:"merges.current,string"`
	MergesCurrentDoc                *int    `json:"merges.current_docs,string"`
	MergesCurrentSize               *string `json:"merges.current_size"`
	MergesTotal                     *int    `json:"merges.total,string"`
	MergesTotalDocs                 *int    `json:"merges.total_docs,string"`
	MergesTotalSize                 *string `json:"merges.total_size"`
	MergesTotalTime                 *string `json:"merges.total_time"`
	RefreshTotal                    *int    `json:"refresh.total,string"`
	RefreshTime                     *string `json:"refresh.time"`
	RefreshExternalTotal            *int    `json:"refresh.external_total,string"`
	RefreshExternalTime             *string `json:"refresh.external_time"`
	RefreshListeners                *int    `json:"refresh.listeners,string"`
	ScriptCompilations              *int    `json:"script.compilations,string"`
	ScriptCacheEvictions            *int    `json:"script.cache_evictions,string"`
	ScriptCompilationLimitTriggered *int    `json:"script.compilation_limit_triggered,string"`
	SearchFetchCurrent              *int    `json:"search.fetch_current,string"`
	SearchFetchTime                 *string `json:"search.fetch_time"`
	SearchFetchTotal                *int    `json:"search.fetch_total,string"`
	SearchOpenContexts              *int    `json:"search.open_contexts,string"`
	SearchQueryCurrent              *int    `json:"search.query_current,string"`
	SearchQueryTime                 *string `json:"search.query_time"`
	SearchQueryTotal                *int    `json:"search.query_total,string"`
	SearchConcurrentQueryCurrent    *int    `json:"search.concurrent_query_current,string"`
	SearchConcurrentQueryTime       *string `json:"search.concurrent_query_time"`
	SearchConcurrentQueryTotal      *int    `json:"search.concurrent_query_total,string"`
	SearchConcurrentAvgSliceCount   *string `json:"search.concurrent_avg_slice_count"`
	SearchScrollCurrent             *int    `json:"search.scroll_current,string"`
	SearchScrollTime                *string `json:"search.scroll_time"`
	SearchScrollTotal               *int    `json:"search.scroll_total,string"`
	SearchPointInTimeCurrent        *int    `json:"search.point_in_time_current,string"`
	SearchPointInTimeTime           *string `json:"search.point_in_time_time"`
	SearchPointInTimeTotal          *int    `json:"search.point_in_time_total,string"`
	SegmentsCount                   *int    `json:"segments.count,string"`
	SegmentsMemory                  *string `json:"segments.memory"`
	SegmentsIndexWriteMemory        *string `json:"segments.index_writer_memory"`
	SegmentsVersionMapMemory        *string `json:"segments.version_map_memory"`
	SegmentsFixedBitsetMemory       *string `json:"segments.fixed_bitset_memory"`
	SuggestCurrent                  *int    `json:"suggest.current,string"`
	SuggestTime                     *string `json:"suggest.time"`
	SuggestTotal                    *int    `json:"suggest.total,string"`
}

func CatNodes(endpoint, nodeName, bytes, time *string) ([]Node, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cat/nodes?format=json&h=name,ip,node.role,node.roles,master,heap.percent,cpu,load_1m,load_5m,load_15m,ram.percent"
	}

	if bytes != nil {
		*endpoint += fmt.Sprintf("&bytes=%s", *bytes)
	}

	if time != nil {
		*endpoint += fmt.Sprintf("&time=%s", *time)
	}

	nodes := make([]Node, 0)

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&nodes).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get nodes: %s", resp.Status())
	}

	if nodeName != nil {
		filtered := make([]Node, 0, len(nodes))

		for _, node := range nodes {
			if strings.Contains(node.Name, *nodeName) {
				filtered = append(filtered, node)
			}
		}

		// If no matches found, return error
		if len(filtered) == 0 {
			return nil, fmt.Errorf("node not found: %s", *nodeName)
		}

		// Replace the original slice with the filtered one
		nodes = filtered
	}

	return nodes, nil
}
