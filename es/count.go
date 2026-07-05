package es

import (
	"context"
	"fmt"
)

// GroupCount represents counts grouped by key.
type GroupCount map[string]int

// IndexGroupCount groups by index name.
type IndexGroupCount map[string]GroupCount

// countAggResponse models the aggregation portion of a _search?size=0 response.
type countAggResponse struct {
	Aggregations struct {
		ByIndex struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
				ByGroup  struct {
					Buckets []struct {
						Key         any    `json:"key"`
						KeyAsString string `json:"key_as_string"`
						DocCount    int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"by_group"`
			} `json:"buckets"`
		} `json:"by_index"`
	} `json:"aggregations"`
}

// buildCountFilters turns --term (field:value) and --exists (field) inputs into
// bool query filters, wrapping fields under a known nested path in a nested query.
func buildCountFilters(termFilters, existsFilters, nestedPaths []string) ([]map[string]any, error) {
	var filters []map[string]any

	for _, term := range termFilters {
		field, value, err := extractFieldAndValue(term)
		if err != nil {
			return nil, err
		}
		clause := map[string]any{"term": map[string]any{field: value}}
		if nestedPath, ok := getNestedPath(field, nestedPaths); ok {
			clause = map[string]any{"nested": map[string]any{"path": nestedPath, "query": clause}}
		}
		filters = append(filters, clause)
	}

	for _, field := range existsFilters {
		clause := map[string]any{"exists": map[string]any{"field": field}}
		if nestedPath, ok := getNestedPath(field, nestedPaths); ok {
			clause = map[string]any{"nested": map[string]any{"path": nestedPath, "query": clause}}
		}
		filters = append(filters, clause)
	}

	return filters, nil
}

// CountDocuments counts documents per index (optionally broken down by a field)
// using a size:0 _search with an _index terms aggregation. Filters from --term
// and --exists narrow the counted set. When groupBy is set, each index bucket is
// sub-aggregated by that field (up to size buckets, default 10).
func CountDocuments(ctx context.Context, index string, termFilters, existsFilters, nestedPaths []string, groupBy string, size int, timeout string, refresh bool) (IndexGroupCount, error) {
	if refresh {
		refreshPath := "_refresh"
		if index != "" {
			refreshPath = fmt.Sprintf("%s/_refresh", index)
		}
		if err := postWithoutBody(ctx, refreshPath, &JsonResponse{}); err != nil {
			return nil, fmt.Errorf("failed to refresh before counting: %w", err)
		}
	}

	filters, err := buildCountFilters(termFilters, existsFilters, nestedPaths)
	if err != nil {
		return nil, err
	}

	// A terms agg on the metadata field _index yields per-index counts; 10000 is
	// well above any realistic index count for a single count invocation.
	byIndex := map[string]any{
		"terms": map[string]any{"field": "_index", "size": 10000},
	}
	if groupBy != "" {
		groupSize := size
		if groupSize <= 0 {
			groupSize = 10
		}
		byIndex["aggs"] = map[string]any{
			"by_group": map[string]any{
				"terms": map[string]any{"field": groupBy, "size": groupSize},
			},
		}
	}

	requestBody := map[string]any{
		"size":  0,
		"query": map[string]any{"bool": map[string]any{"filter": filters}},
		"aggs":  map[string]any{"by_index": byIndex},
	}
	if timeout != "" {
		requestBody["timeout"] = timeout
	}

	endpoint := "_search"
	if index != "" {
		endpoint = fmt.Sprintf("%s/_search", index)
	}

	var resp countAggResponse
	if err := postJSONResponseWithBody(ctx, endpoint, &resp, requestBody); err != nil {
		return nil, err
	}

	result := IndexGroupCount{}
	for _, idxBucket := range resp.Aggregations.ByIndex.Buckets {
		groups := GroupCount{}
		if groupBy == "" {
			groups[""] = idxBucket.DocCount
		} else {
			for _, gb := range idxBucket.ByGroup.Buckets {
				key := gb.KeyAsString
				if key == "" {
					key = fmt.Sprintf("%v", gb.Key)
				}
				groups[key] = gb.DocCount
			}
		}
		result[idxBucket.Key] = groups
	}

	return result, nil
}
