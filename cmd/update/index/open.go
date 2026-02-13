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
	flagOpenIndices string
	flagOpenDryRun  bool
)

var OpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open one or more closed indices",
	Long: utils.Trim(`
		Opens a closed index. A closed index is blocked for read/write operations and does not allow
		all operations that opened indices allow. It is not possible to index documents or to search
		for documents in a closed index.
	`),
	Example: utils.TrimAndIndent(`
		# Open a single index
		esctl update index open --indices my-index

		# Open multiple indices
		esctl update index open --indices index1,index2,index3

		# Dry run - show what would be opened
		esctl update index open --indices my-index --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleOpenIndex(ctx)
	},
}

func init() {
	OpenCmd.Flags().StringVar(&flagOpenIndices, "indices", "", "Comma-separated list of indices to open")
	OpenCmd.Flags().BoolVar(&flagOpenDryRun, "dry-run", false, "Show what would be opened without actually opening")
	_ = OpenCmd.MarkFlagRequired("indices")
}

func handleOpenIndex(ctx context.Context) error {
	indices, err := utils.ParseIndexPatternsCSV(flagOpenIndices, true)
	if err != nil {
		return err
	}

	// Validate each index name
	for _, idx := range indices {
		if err := validation.ValidateIndexPattern(idx); err != nil {
			return fmt.Errorf("invalid index pattern %s: %w", idx, err)
		}
	}

	if flagOpenDryRun {
		fmt.Printf("DRY RUN: Would open the following indices:\n")
		for _, idx := range indices {
			fmt.Printf("  - %s\n", idx)
		}
		return nil
	}

	idxClient := index.NewIndex()
	resp, err := idxClient.Open(ctx, indices)
	if err != nil {
		return fmt.Errorf("failed to open indices: %w", err)
	}

	if resp.Acknowledged {
		fmt.Printf("Successfully opened %d index(es)\n", len(indices))
	}

	return output.Render(resp)
}
