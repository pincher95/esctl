package alias

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
)

func HandleAliasAdd(ctx context.Context, alias, indicesCSV, filterJSON, routing string) error {
	if indicesCSV == "" {
		return fmt.Errorf("indices must be specified")
	}

	if err := validation.ValidateAliasName(alias); err != nil {
		return err
	}

	indices := strings.Split(indicesCSV, ",")
	trimmed := make([]string, 0, len(indices))
	for _, idx := range indices {
		clean := strings.TrimSpace(idx)
		if clean == "" {
			continue
		}
		if err := validation.ValidateIndexPattern(clean); err != nil {
			return err
		}
		trimmed = append(trimmed, clean)
	}
	if len(trimmed) == 0 {
		return fmt.Errorf("indices must be specified")
	}

	var filter map[string]interface{}
	if filterJSON != "" {
		if err := json.Unmarshal([]byte(filterJSON), &filter); err != nil {
			return fmt.Errorf("invalid filter JSON: %w", err)
		}
	}

	if err := index.AddAlias(ctx, trimmed, alias, filter, routing); err != nil {
		return fmt.Errorf("failed to add alias: %w", err)
	}

	fmt.Printf("Successfully added alias '%s' to indices: %s\n", alias, strings.Join(trimmed, ", "))
	return nil
}
