package es

import (
	"context"
	"fmt"
	"strings"
)

type JsonResponse map[string]any

func extractFieldAndValue(term string) (string, string, error) {
	parts := strings.SplitN(term, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid term format: %s", term)
	}
	return parts[0], parts[1], nil
}

func SearchDocuments(
	ctx context.Context,
	index string,
	ids []string,
	terms []string,
	from int,
	size int,
	nestedPaths []string,
	sortFields []string,
) (JsonResponse, error) {
	var filters []map[string]any

	for _, term := range terms {
		field, value, err := extractFieldAndValue(term)
		if err != nil {
			return nil, err
		}
		nestedPath, isNestedPath := getNestedPath(field, nestedPaths)
		if isNestedPath {
			termFilter := map[string]any{
				"nested": map[string]any{
					"path": nestedPath,
					"query": map[string]any{
						"term": map[string]any{
							field: value,
						},
					},
				},
			}
			filters = append(filters, termFilter)
		} else {
			termFilter := map[string]any{
				"term": map[string]any{
					field: value,
				},
			}
			filters = append(filters, termFilter)
		}
	}

	if len(ids) > 0 {
		idsFilter := map[string]any{
			"ids": map[string]any{
				"values": ids,
			},
		}
		filters = append(filters, idsFilter)
	}

	query := map[string]any{
		"bool": map[string]any{
			"filter": filters,
		},
	}

	requestBody := map[string]any{
		"from":  from,
		"size":  max(size, len(ids)),
		"query": query,
	}

	if len(sortFields) > 0 {
		sorts := make([]map[string]string, len(sortFields))
		for i, sortField := range sortFields {
			field, order, err := extractFieldAndValue(sortField)
			if err != nil {
				return nil, err
			}
			sorts[i] = map[string]string{field: order}
		}
		requestBody["sort"] = sorts
	}

	endpoint := fmt.Sprintf("%s/_search", index)
	var response JsonResponse
	err := postJSONResponseWithBody(ctx, endpoint, &response, requestBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}
