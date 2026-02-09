package alias

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
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

	trimmed, err := utils.ParseIndexPatternsCSV(indicesCSV, true)
	if err != nil {
		return err
	}

	var filter map[string]any
	if filterJSON != "" {
		filter, err = utils.ParseJSONMap(filterJSON, "invalid filter JSON")
		if err != nil {
			return err
		}
	}

	if err := index.AddAlias(ctx, trimmed, alias, filter, routing); err != nil {
		return fmt.Errorf("failed to add alias: %w", err)
	}

	fmt.Printf("Successfully added alias '%s' to indices: %s\n", alias, strings.Join(trimmed, ", "))
	return nil
}
