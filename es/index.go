package es

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/cat"
)

type IndexMappings struct {
	Mappings any `json:"mappings"`
}

type IndexSettings struct {
	Settings any `json:"settings"`
}

type MappingsResponse map[string]IndexMappings
type SettingsResponse map[string]IndexSettings

type IndexDetails struct {
	Settings any `json:"settings,omitempty"`
	Mappings any `json:"mappings,omitempty"`
}

type IndexDetailsResponse map[string]IndexDetails

func GetIndexDetails(index string, shouldGetMappings, shouldGetSettings bool) (IndexDetailsResponse, error) {
	var mappingsResponse MappingsResponse
	var settingsResponse SettingsResponse

	if shouldGetMappings {
		mappingsEndpoint := fmt.Sprintf("%s/_mappings", index)
		if err := getJSONResponse(mappingsEndpoint, &mappingsResponse); err != nil {
			return nil, fmt.Errorf("failed to get index mappings: %w", err)
		}
	}

	if shouldGetSettings {
		settingsEndpoint := fmt.Sprintf("%s/_settings", index)
		if err := getJSONResponse(settingsEndpoint, &settingsResponse); err != nil {
			return nil, fmt.Errorf("failed to get index settings: %w", err)
		}
	}

	// Merge settings and mappings.
	merged := make(IndexDetailsResponse)
	for index, indexMappings := range mappingsResponse {
		details := IndexDetails{
			Mappings: indexMappings.Mappings,
		}

		if indexSettings, ok := settingsResponse[index]; ok {
			details.Settings = indexSettings.Settings
		}

		merged[index] = details
	}

	// Process settings for indices that don't exist in mappings.
	for index, indexSettings := range settingsResponse {
		if _, exists := merged[index]; !exists {
			merged[index] = IndexDetails{
				Settings: indexSettings.Settings, // Only settings are added.
			}
		}
	}

	return merged, nil
}

type AliasResponse map[string]AliasDetail

type AliasDetail struct {
	Aliases map[string]any `json:"aliases"`
}

func GetAliases(index string) (map[string]string, error) {
	if index == "" {
		index = "_all"
	}

	var aliasResp AliasResponse
	if err := getJSONResponse(index+"/_alias", &aliasResp); err != nil {
		return nil, err
	}

	aliases := make(map[string]string)
	for index, detail := range aliasResp {
		for alias := range detail.Aliases {
			aliases[alias] = index
		}
	}

	return aliases, nil
}

type CountResponse struct {
	Count        int            `json:"count"`
	Aggregations map[string]any `json:"aggregations,omitempty"`
}

type GroupCount map[string]int
type IndexGroupCount map[string]GroupCount

func buildFilterQueries(termFilters, existsFilters []string, nestedPaths []string) []map[string]any {
	filterQueries := make([]map[string]any, 0)
	nestedGroups := make(map[string][]map[string]any)

	groupFilter := func(filter string, isTerm bool) {
		var field, value string
		var filterQuery map[string]any

		if isTerm {
			parts := strings.SplitN(filter, ":", 2)
			if len(parts) != 2 {
				return
			}
			field, value = parts[0], parts[1]
			filterQuery = map[string]any{
				"term": map[string]any{
					field: value,
				},
			}
		} else {
			field = filter
			filterQuery = map[string]any{
				"exists": map[string]any{
					"field": field,
				},
			}
		}

		nestedPath, isNestedPath := getNestedPath(field, nestedPaths)
		if isNestedPath {
			nestedGroups[nestedPath] = append(nestedGroups[nestedPath], filterQuery)
		} else {
			filterQueries = append(filterQueries, filterQuery)
		}
	}

	for _, filter := range termFilters {
		groupFilter(filter, true)
	}
	for _, filter := range existsFilters {
		groupFilter(filter, false)
	}

	for path, groupedFilters := range nestedGroups {
		filterQueries = append(filterQueries, map[string]any{
			"nested": map[string]any{
				"path": path,
				"query": map[string]any{
					"bool": map[string]any{
						"must": groupedFilters,
					},
				},
			},
		})
	}

	return filterQueries
}

func countDocumentsOfIndex(index string, termFilters, existsFilters, nestedPaths []string) (int, error) {
	endpoint := index + "/_count"
	query := map[string]any{
		"match_all": map[string]any{},
	}

	if len(termFilters) > 0 || len(existsFilters) > 0 {
		query = map[string]any{
			"bool": map[string]any{
				"must": buildFilterQueries(termFilters, existsFilters, nestedPaths),
			},
		}
	}

	body := map[string]any{
		"query": query,
	}

	var response CountResponse
	if err := postJSONResponseWithBody(endpoint, &response, body); err != nil {
		return 0, err
	}

	return response.Count, nil
}

type RefreshResponse map[string]any

func RefreshIndices(target string) error {
	endpoint := target + "/_refresh"
	var response RefreshResponse
	return postWithoutBody(endpoint, &response)
}

func groupDocumentsOfIndex(
	index string,
	termFilters []string,
	existsFilters []string,
	nestedPaths []string,
	groupBy string,
	size int,
	timeout string,
) (GroupCount, error) {
	endpoint := index + "/_search"
	query := map[string]any{
		"match_all": map[string]any{},
	}

	if len(termFilters) > 0 || len(existsFilters) > 0 {
		query = map[string]any{
			"bool": map[string]any{
				"must": buildFilterQueries(termFilters, existsFilters, nestedPaths),
			},
		}
	}

	if size <= 0 {
		size = 50
	}

	if timeout == "" {
		timeout = "1s"
	}

	nestedPath, isNestedPath := getNestedPath(groupBy, nestedPaths)
	aggregations := make(map[string]any)

	if isNestedPath {
		aggregations["group_by_nested"] = map[string]any{
			"nested": map[string]any{
				"path": nestedPath,
			},
			"aggs": map[string]any{
				"group_by": map[string]any{
					"terms": map[string]any{
						"field": groupBy,
						"size":  size,
					},
				},
			},
		}
	} else {
		aggregations["group_by"] = map[string]any{
			"terms": map[string]any{
				"field": groupBy,
				"size":  size,
			},
		}
	}

	body := map[string]any{
		"query":   query,
		"aggs":    aggregations,
		"timeout": timeout,
	}

	var response CountResponse
	if err := postJSONResponseWithBody(endpoint, &response, body); err != nil {
		return nil, err
	}

	groupCount := make(GroupCount)
	var buckets map[string]any

	if isNestedPath {
		if nestedAgg, ok := response.Aggregations["group_by_nested"].(map[string]any); ok {
			buckets, _ = nestedAgg["group_by"].(map[string]any)
		}
	} else {
		buckets, _ = response.Aggregations["group_by"].(map[string]any)
	}

	if terms, ok := buckets["buckets"].([]any); ok {
		for _, term := range terms {
			if termData, ok := term.(map[string]any); ok {
				key := fmt.Sprint(termData["key"])
				count := int(termData["doc_count"].(float64))
				groupCount[key] = count
			}
		}
	}

	return groupCount, nil
}

func CountDocuments(
	index string,
	termFilters []string,
	existsFilters []string,
	nestedPaths []string,
	groupBy string,
	size int,
	timeout string,
	refresh bool,
) (map[string]GroupCount, error) {
	if refresh {
		err := RefreshIndices(index)
		if err != nil {
			return nil, err
		}
	}

	client := cat.NewCat()
	indices, err := client.CatIndices(nil, &index, nil)
	if err != nil {
		return nil, err
	}

	indexCounts := make(map[string]GroupCount)
	for _, index := range *indices {
		var groupCount GroupCount
		if groupBy == "" {
			count, err := countDocumentsOfIndex(index.Index, termFilters, existsFilters, nestedPaths)
			if err != nil {
				return nil, err
			}
			groupCount = map[string]int{"": count}
		} else {
			groupCount, err = groupDocumentsOfIndex(index.Index, termFilters, existsFilters, nestedPaths, groupBy, size, timeout)
			if err != nil {
				return nil, err
			}
		}

		indexCounts[index.Index] = groupCount
	}

	return indexCounts, nil
}
