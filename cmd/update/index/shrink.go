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
	flagShrinkSource   string
	flagShrinkTarget   string
	flagShrinkShards   int
	flagShrinkSettings string
	flagShrinkDryRun   bool
)

var ShrinkCmd = &cobra.Command{
	Use:   "shrink",
	Short: "Shrink an index into a new index with fewer primary shards",
	Long: utils.Trim(`
		Shrinks an existing index into a new index with fewer primary shards. The source index must
		be marked as read-only and have green health. All primary shards must be relocated to the
		same node before shrinking. The target shard count must be a factor of the source shard count.
	`),
	Example: utils.TrimAndIndent(`
		# Shrink an index from 6 shards to 1 shard
		esctl update index shrink --source my-index --target my-index-shrink --shards 1

		# Shrink with custom settings
		esctl update index shrink --source my-index --target my-index-shrink --shards 1 --settings '{"index.number_of_replicas":2}'

		# Dry run - show what would be shrunk
		esctl update index shrink --source my-index --target my-index-shrink --shards 1 --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleShrinkIndex(ctx)
	},
}

func init() {
	ShrinkCmd.Flags().StringVar(&flagShrinkSource, "source", "", "Source index name")
	ShrinkCmd.Flags().StringVar(&flagShrinkTarget, "target", "", "Target index name")
	ShrinkCmd.Flags().IntVar(&flagShrinkShards, "shards", 0, "Number of primary shards for target index")
	ShrinkCmd.Flags().StringVar(&flagShrinkSettings, "settings", "", "Additional settings as JSON")
	ShrinkCmd.Flags().BoolVar(&flagShrinkDryRun, "dry-run", false, "Show what would be shrunk without actually shrinking")
	_ = ShrinkCmd.MarkFlagRequired("source")
	_ = ShrinkCmd.MarkFlagRequired("target")
	_ = ShrinkCmd.MarkFlagRequired("shards")
}

func handleShrinkIndex(ctx context.Context) error {
	// Validate source and target
	if err := validation.ValidateIndexName(flagShrinkSource); err != nil {
		return fmt.Errorf("invalid source index name: %w", err)
	}
	if err := validation.ValidateIndexName(flagShrinkTarget); err != nil {
		return fmt.Errorf("invalid target index name: %w", err)
	}

	// Validate shard count
	if err := validation.ValidateShardCount(flagShrinkShards); err != nil {
		return fmt.Errorf("invalid shard count: %w", err)
	}

	// Parse settings if provided
	var settings map[string]any
	if flagShrinkSettings != "" {
		var err error
		settings, err = utils.ParseJSONMap(flagShrinkSettings, "invalid settings JSON")
		if err != nil {
			return err
		}
	}

	if flagShrinkDryRun {
		fmt.Printf("DRY RUN: Would shrink index:\n")
		fmt.Printf("  Source: %s\n", flagShrinkSource)
		fmt.Printf("  Target: %s\n", flagShrinkTarget)
		fmt.Printf("  Target shards: %d\n", flagShrinkShards)
		if settings != nil {
			fmt.Printf("  Settings: %v\n", settings)
		}
		return nil
	}

	idxClient := index.NewIndex()
	resp, err := idxClient.Shrink(ctx, flagShrinkSource, flagShrinkTarget, flagShrinkShards, settings)
	if err != nil {
		return fmt.Errorf("failed to shrink index: %w", err)
	}

	if resp.Acknowledged {
		fmt.Printf("Successfully shrunk %s to %s with %d shard(s)\n", flagShrinkSource, flagShrinkTarget, flagShrinkShards)
	}

	return output.Render(resp)
}
