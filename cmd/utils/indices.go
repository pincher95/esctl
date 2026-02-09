package utils

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/internal/validation"
)

// ParseIndexPatternsCSV splits comma-delimited index patterns and validates each.
func ParseIndexPatternsCSV(indicesCSV string, required bool) ([]string, error) {
	if indicesCSV == "" {
		if required {
			return nil, fmt.Errorf("indices must be specified")
		}
		return nil, nil
	}

	parts := strings.Split(indicesCSV, ",")
	trimmed := make([]string, 0, len(parts))
	for _, idx := range parts {
		clean := strings.TrimSpace(idx)
		if clean == "" {
			continue
		}
		if err := validation.ValidateIndexPattern(clean); err != nil {
			return nil, err
		}
		trimmed = append(trimmed, clean)
	}
	if required && len(trimmed) == 0 {
		return nil, fmt.Errorf("indices must be specified")
	}
	return trimmed, nil
}

// ValidateIndexPatternsCSV validates each index pattern in a CSV string.
func ValidateIndexPatternsCSV(indicesCSV string) error {
	_, err := ParseIndexPatternsCSV(indicesCSV, false)
	return err
}
