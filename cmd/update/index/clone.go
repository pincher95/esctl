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
	flagCloneSource   string
	flagCloneTarget   string
	flagCloneSettings string
	flagCloneDryRun   bool
)

var CloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone an index to a new index",
	Long: utils.Trim(`
		Clones an existing index into a new index. The source index must be marked as read-only
		and have green health to be cloned. The clone operation creates a new index with the same
		data and settings as the source index.
	`),
	Example: utils.TrimAndIndent(`
		# Clone an index
		esctl update index clone --source my-index --target my-index-clone

		# Clone with custom settings
		esctl update index clone --source my-index --target my-index-clone --settings '{"index.number_of_replicas":2}'

		# Dry run - show what would be cloned
		esctl update index clone --source my-index --target my-index-clone --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleCloneIndex(ctx)
	},
}

func init() {
	CloneCmd.Flags().StringVar(&flagCloneSource, "source", "", "Source index name")
	CloneCmd.Flags().StringVar(&flagCloneTarget, "target", "", "Target index name")
	CloneCmd.Flags().StringVar(&flagCloneSettings, "settings", "", "Additional settings as JSON")
	CloneCmd.Flags().BoolVar(&flagCloneDryRun, "dry-run", false, "Show what would be cloned without actually cloning")
	_ = CloneCmd.MarkFlagRequired("source")
	_ = CloneCmd.MarkFlagRequired("target")
}

func handleCloneIndex(ctx context.Context) error {
	// Validate source and target
	if err := validation.ValidateIndexName(flagCloneSource); err != nil {
		return fmt.Errorf("invalid source index name: %w", err)
	}
	if err := validation.ValidateIndexName(flagCloneTarget); err != nil {
		return fmt.Errorf("invalid target index name: %w", err)
	}

	// Parse settings if provided
	var settings map[string]any
	if flagCloneSettings != "" {
		var err error
		settings, err = utils.ParseJSONMap(flagCloneSettings, "invalid settings JSON")
		if err != nil {
			return err
		}
	}

	if flagCloneDryRun {
		fmt.Printf("DRY RUN: Would clone index:\n")
		fmt.Printf("  Source: %s\n", flagCloneSource)
		fmt.Printf("  Target: %s\n", flagCloneTarget)
		if settings != nil {
			fmt.Printf("  Settings: %v\n", settings)
		}
		return nil
	}

	idxClient := index.NewIndex()
	resp, err := idxClient.Clone(ctx, flagCloneSource, flagCloneTarget, settings)
	if err != nil {
		return fmt.Errorf("failed to clone index: %w", err)
	}

	if resp.Acknowledged {
		fmt.Printf("Successfully cloned %s to %s\n", flagCloneSource, flagCloneTarget)
	}

	return output.Render(resp)
}
