package alias

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move <alias> --from=<from-index> --to=<to-index>",
	Short: "Atomically move an alias from one index to another",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Move alias from old index to new index
	esctl alias move logs-current --from=logs-2023 --to=logs-2024

	# Move alias during zero-downtime reindexing
	esctl alias move my-app --from=my-app-v1 --to=my-app-v2
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleAliasMove(ctx, args[0])
	},
}

var (
	moveFrom string
	moveTo   string
)

func init() {
	moveCmd.Flags().StringVar(&moveFrom, "from", "", "Index to move the alias from")
	moveCmd.Flags().StringVar(&moveTo, "to", "", "Index to move the alias to")
	moveCmd.MarkFlagRequired("from")
	moveCmd.MarkFlagRequired("to")
}

func handleAliasMove(ctx context.Context, alias string) error {
	if moveFrom == "" || moveTo == "" {
		return fmt.Errorf("both --from and --to indices must be specified")
	}

	err := index.MoveAlias(ctx, moveFrom, moveTo, alias)
	if err != nil {
		return fmt.Errorf("failed to move alias: %w", err)
	}

	fmt.Printf("Successfully moved alias '%s' from '%s' to '%s'\n", alias, moveFrom, moveTo)
	return nil
}
