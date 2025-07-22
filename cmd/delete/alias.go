package delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/delete"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias <alias> --indices=<index1,index2,...>",
	Short: "Remove an alias from indices",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Remove alias from a single index
	esctl delete alias my-alias --indices=my-index

	# Remove alias from multiple indices
	esctl delete alias logs-current --indices="logs-2023,logs-2024"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleAliasDelete(ctx, args[0])
	},
}

var (
	aliasIndices string
)

func init() {
	aliasCmd.Flags().StringVar(&aliasIndices, "indices", "", "Comma-separated list of indices to remove the alias from")
	aliasCmd.MarkFlagRequired("indices")
}

func handleAliasDelete(ctx context.Context, alias string) error {
	if aliasIndices == "" {
		return fmt.Errorf("indices must be specified")
	}

	indices := strings.Split(aliasIndices, ",")
	for i, idx := range indices {
		indices[i] = strings.TrimSpace(idx)
	}

	aliases := []string{alias}

	err := delete.DeleteAlias(ctx, indices, aliases)
	if err != nil {
		return fmt.Errorf("failed to delete alias: %w", err)
	}

	fmt.Printf("Successfully removed alias '%s' from indices: %s\n", alias, strings.Join(indices, ", "))
	return nil
}
