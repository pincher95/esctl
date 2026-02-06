package alias

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
)

func HandleAliasRemove(ctx context.Context, alias, indicesCSV string) error {
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

	if err := index.RemoveAlias(ctx, trimmed, alias); err != nil {
		return fmt.Errorf("failed to remove alias: %w", err)
	}

	fmt.Printf("Successfully removed alias '%s' from indices: %s\n", alias, strings.Join(trimmed, ", "))
	return nil
}
