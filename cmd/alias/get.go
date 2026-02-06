package alias

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/pincher95/esctl/output"
)

func HandleAliasGet(ctx context.Context, aliasName string) error {
	if err := validation.ValidateAliasName(aliasName); err != nil {
		return err
	}

	aliases, err := index.GetAlias(ctx, nil, []string{aliasName})
	if err != nil {
		return fmt.Errorf("failed to get alias: %w", err)
	}

	if len(aliases) == 0 {
		return fmt.Errorf("alias not found: %s", aliasName)
	}

	return output.Render(aliases)
}
