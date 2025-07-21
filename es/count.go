package es

import (
	"context"
)

// GroupCount represents counts grouped by key.
type GroupCount map[string]int

// IndexGroupCount groups by index name.
type IndexGroupCount map[string]GroupCount

// CountDocuments is a stub for now to keep CLI compiling; returns zero counts.
func CountDocuments(ctx context.Context, index string, termFilters, existsFilters, nestedPaths []string, groupBy string, size int, timeout string, refresh bool) (IndexGroupCount, error) {
	return IndexGroupCount{}, nil
}
