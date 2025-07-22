package alias

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	listIndexPattern string
	listAliasPattern string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases",
	Example: utils.TrimAndIndent(`
	# List all aliases
	esctl alias list

	# List aliases for specific index pattern
	esctl alias list --index-pattern="logs-*"

	# List specific alias pattern
	esctl alias list --alias-pattern="*-current"

	# List aliases matching both patterns
	esctl alias list --index-pattern="logs-*" --alias-pattern="*-current"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleAliasList(ctx)
	},
}

func init() {
	listCmd.Flags().StringVar(&listIndexPattern, "index-pattern", "", "Pattern to match index names")
	listCmd.Flags().StringVar(&listAliasPattern, "alias-pattern", "", "Pattern to match alias names")
}

func handleAliasList(ctx context.Context) error {
	result, err := index.ListAliases(ctx, listIndexPattern, listAliasPattern)
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(result) == 0 {
		fmt.Println("No aliases found")
		return nil
	}

	output.PrintJson(result)
	return nil
}
