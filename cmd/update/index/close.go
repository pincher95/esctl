package index

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagCloseIndices string
	flagCloseDryRun  bool
)

var CloseCmd = &cobra.Command{
	Use:   "close",
	Short: "Close one or more open indices",
	Long: utils.Trim(`
		Closes an index. A closed index is blocked for read/write operations and does not allow
		all operations that opened indices allow. Closing indices can be useful to reduce the overhead
		on the cluster when indices are not actively being used.
	`),
	Example: utils.TrimAndIndent(`
		# Close a single index
		esctl update index close --indices my-index

		# Close multiple indices
		esctl update index close --indices index1,index2,index3

		# Dry run - show what would be closed
		esctl update index close --indices my-index --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleCloseIndex(ctx)
	},
}

func init() {
	CloseCmd.Flags().StringVar(&flagCloseIndices, "indices", "", "Comma-separated list of indices to close")
	CloseCmd.Flags().BoolVar(&flagCloseDryRun, "dry-run", false, "Show what would be closed without actually closing")
	_ = CloseCmd.MarkFlagRequired("indices")
}

func handleCloseIndex(ctx context.Context) error {
	indices, err := utils.ParseIndexPatternsCSV(flagCloseIndices, true)
	if err != nil {
		return err
	}

	// Validate each index name
	for _, idx := range indices {
		if err := validation.ValidateIndexPattern(idx); err != nil {
			return fmt.Errorf("invalid index pattern %s: %w", idx, err)
		}
	}

	if flagCloseDryRun {
		fmt.Printf("DRY RUN: Would close the following indices:\n")
		for _, idx := range indices {
			fmt.Printf("  - %s\n", idx)
		}
		return nil
	}

	idxClient := index.NewIndex()
	resp, err := idxClient.Close(ctx, indices)
	if err != nil {
		return fmt.Errorf("failed to close indices: %w", err)
	}

	if resp.Acknowledged {
		fmt.Printf("Successfully closed %d index(es)\n", len(indices))
	}

	return output.Render(resp)
}
