package alias

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/spf13/cobra"
)

var (
	removeIndices string
)

var removeCmd = &cobra.Command{
	Use:   "remove <alias> --indices=<index1,index2,...>",
	Short: "Remove an alias from one or more indices",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Remove alias from a single index
	esctl alias remove my-alias --indices=my-index

	# Remove alias from multiple indices
	esctl alias remove logs-current --indices="logs-2023,logs-2024"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleAliasRemove(ctx, args[0])
	},
}

func init() {
	removeCmd.Flags().StringVar(&removeIndices, "indices", "", "Comma-separated list of indices to remove the alias from")
	removeCmd.MarkFlagRequired("indices")
}

func handleAliasRemove(ctx context.Context, alias string) error {
	if removeIndices == "" {
		return fmt.Errorf("indices must be specified")
	}

	indices := strings.Split(removeIndices, ",")
	for i, idx := range indices {
		indices[i] = strings.TrimSpace(idx)
	}

	err := index.RemoveAlias(ctx, indices, alias)
	if err != nil {
		return fmt.Errorf("failed to remove alias: %w", err)
	}

	fmt.Printf("Successfully removed alias '%s' from indices: %s\n", alias, strings.Join(indices, ", "))
	return nil
}
