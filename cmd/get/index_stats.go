package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagStatsIndices string
	flagStatsMetric  string
)

var getIndexStatsCmd = &cobra.Command{
	Use:   "index-stats",
	Short: "Get detailed statistics for indices",
	Long: utils.Trim(`
		Retrieves comprehensive statistics for one or more indices, including document counts,
		storage sizes, indexing rates, search performance, merge statistics, and more.
	`),
	Example: utils.TrimAndIndent(`
		# Get stats for all indices
		esctl get index-stats

		# Get stats for specific indices
		esctl get index-stats --indices my-index,other-index

		# Get specific metric only (docs, store, indexing, search, etc.)
		esctl get index-stats --indices my-index --metric docs

		# Watch mode for monitoring
		esctl get index-stats --indices my-index --watch

		# JSON output
		esctl get index-stats --indices my-index -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idxClient := index.NewIndex()

		if !flagRefresh {
			return handleIndexStatsLogic(ctx, idxClient)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleIndexStatsLogic(ctx, idxClient)
		})
	},
}

func init() {
	getIndexStatsCmd.Flags().StringVar(&flagStatsIndices, "indices", "", "Comma-separated list of indices")
	getIndexStatsCmd.Flags().StringVar(&flagStatsMetric, "metric", "", "Specific metric to retrieve (docs, store, indexing, search, etc.)")
}

func handleIndexStatsLogic(ctx context.Context, idxClient index.Index) error {
	var indices []string
	if flagStatsIndices != "" {
		var err error
		indices, err = utils.ParseIndexPatternsCSV(flagStatsIndices, false)
		if err != nil {
			return err
		}
	}

	stats, err := idxClient.GetIndexStats(ctx, indices, flagStatsMetric)
	if err != nil {
		return fmt.Errorf("failed to retrieve index stats: %w", err)
	}

	return output.Render(stats)
}
