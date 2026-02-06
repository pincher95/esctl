package alias

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
)

func HandleAliasMove(ctx context.Context, alias, fromIndex, toIndex string) error {
	if fromIndex == "" || toIndex == "" {
		return fmt.Errorf("both --from and --to indices must be specified")
	}

	if err := validation.ValidateAliasName(alias); err != nil {
		return err
	}
	if err := validation.ValidateIndexPattern(fromIndex); err != nil {
		return err
	}
	if err := validation.ValidateIndexPattern(toIndex); err != nil {
		return err
	}

	if err := index.MoveAlias(ctx, fromIndex, toIndex, alias); err != nil {
		return fmt.Errorf("failed to move alias: %w", err)
	}

	fmt.Printf("Successfully moved alias '%s' from '%s' to '%s'\n", alias, fromIndex, toIndex)
	return nil
}
