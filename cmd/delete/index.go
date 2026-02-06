package delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/delete"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/spf13/cobra"
)

var (
	indexIgnoreUnavailable bool
	indexAllowNoIndices    bool
	indexExpandWildcards   string
	indexForce             bool
)

var indexCmd = &cobra.Command{
	Use:   "index <index1> [index2] [index3]...",
	Short: "Delete one or more indices",
	Args:  cobra.MinimumNArgs(1),
	Example: utils.TrimAndIndent(`
	# Delete a single index
	esctl delete index my-index

	# Delete multiple indices
	esctl delete index index1 index2 index3

	# Delete indices with wildcard pattern
	esctl delete index "logs-*" --expand-wildcards=open

	# Force delete ignoring unavailable indices
	esctl delete index "missing-*" --ignore-unavailable --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleIndexDelete(ctx, args)
	},
}

func init() {
	indexCmd.Flags().BoolVar(&indexIgnoreUnavailable, "ignore-unavailable", false, "Ignore unavailable indices")
	indexCmd.Flags().BoolVar(&indexAllowNoIndices, "allow-no-indices", true, "Allow operations on no indices")
	indexCmd.Flags().StringVar(&indexExpandWildcards, "expand-wildcards", "open", "Expand wildcard expressions (all, open, closed, hidden, none)")
	indexCmd.Flags().BoolVar(&indexForce, "force", false, "Force deletion without confirmation")
}

func handleIndexDelete(ctx context.Context, indices []string) error {
	for _, idx := range indices {
		if err := validation.ValidateIndexPattern(idx); err != nil {
			return err
		}
	}

	if !indexForce {
		fmt.Printf("WARNING: This will permanently delete the following indices: %s\n", strings.Join(indices, ", "))
		fmt.Print("Are you sure you want to continue? (y/N): ")
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" && strings.ToLower(confirmation) != "yes" {
			fmt.Println("Operation cancelled")
			return nil
		}
	}

	err := delete.DeleteIndex(ctx, indices, indexIgnoreUnavailable, indexAllowNoIndices, indexExpandWildcards)
	if err != nil {
		return fmt.Errorf("failed to delete indices: %w", err)
	}

	fmt.Printf("Successfully deleted indices: %s\n", strings.Join(indices, ", "))
	return nil
}
