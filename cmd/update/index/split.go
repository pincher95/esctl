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
	flagSplitSource   string
	flagSplitTarget   string
	flagSplitShards   int
	flagSplitSettings string
	flagSplitDryRun   bool
)

var SplitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split an index into a new index with more primary shards",
	Long: utils.Trim(`
		Splits an existing index into a new index with more primary shards. The source index must
		be marked as read-only and have green health. The number of target shards must be a multiple
		of the source shard count.
	`),
	Example: utils.TrimAndIndent(`
		# Split an index from 3 shards to 6 shards
		esctl update index split --source my-index --target my-index-split --shards 6

		# Split with custom settings
		esctl update index split --source my-index --target my-index-split --shards 6 --settings '{"index.number_of_replicas":2}'

		# Dry run - show what would be split
		esctl update index split --source my-index --target my-index-split --shards 6 --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSplitIndex(ctx)
	},
}

func init() {
	SplitCmd.Flags().StringVar(&flagSplitSource, "source", "", "Source index name")
	SplitCmd.Flags().StringVar(&flagSplitTarget, "target", "", "Target index name")
	SplitCmd.Flags().IntVar(&flagSplitShards, "shards", 0, "Number of primary shards for target index")
	SplitCmd.Flags().StringVar(&flagSplitSettings, "settings", "", "Additional settings as JSON")
	SplitCmd.Flags().BoolVar(&flagSplitDryRun, "dry-run", false, "Show what would be split without actually splitting")
	_ = SplitCmd.MarkFlagRequired("source")
	_ = SplitCmd.MarkFlagRequired("target")
	_ = SplitCmd.MarkFlagRequired("shards")
}

func handleSplitIndex(ctx context.Context) error {
	// Validate source and target
	if err := validation.ValidateIndexName(flagSplitSource); err != nil {
		return fmt.Errorf("invalid source index name: %w", err)
	}
	if err := validation.ValidateIndexName(flagSplitTarget); err != nil {
		return fmt.Errorf("invalid target index name: %w", err)
	}

	// Validate shard count
	if err := validation.ValidateShardCount(flagSplitShards); err != nil {
		return fmt.Errorf("invalid shard count: %w", err)
	}

	// Parse settings if provided
	var settings map[string]any
	if flagSplitSettings != "" {
		var err error
		settings, err = utils.ParseJSONMap(flagSplitSettings, "invalid settings JSON")
		if err != nil {
			return err
		}
	}

	if flagSplitDryRun {
		fmt.Printf("DRY RUN: Would split index:\n")
		fmt.Printf("  Source: %s\n", flagSplitSource)
		fmt.Printf("  Target: %s\n", flagSplitTarget)
		fmt.Printf("  Target shards: %d\n", flagSplitShards)
		if settings != nil {
			fmt.Printf("  Settings: %v\n", settings)
		}
		return nil
	}

	idxClient := index.NewIndex()
	resp, err := idxClient.Split(ctx, flagSplitSource, flagSplitTarget, flagSplitShards, settings)
	if err != nil {
		return fmt.Errorf("failed to split index: %w", err)
	}

	if resp.Acknowledged {
		fmt.Printf("Successfully split %s to %s with %d shards\n", flagSplitSource, flagSplitTarget, flagSplitShards)
	}

	return output.Render(resp)
}
