package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flagSegmentsIndices string

var getSegmentsCmd = &cobra.Command{
	Use:   "segments",
	Short: "Get segment information for indices",
	Long: utils.Trim(`
		Retrieves low-level segment information for indices. Segments are the internal data structures
		where Lucene stores index data. This command is useful for understanding merge behavior,
		memory usage, and troubleshooting performance issues.
	`),
	Example: utils.TrimAndIndent(`
		# Get segments for all indices
		esctl get segments

		# Get segments for specific indices
		esctl get segments --indices my-index,other-index

		# Watch mode for monitoring segment merges
		esctl get segments --indices my-index --watch --interval 10s

		# JSON output
		esctl get segments --indices my-index -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idxClient := index.NewIndex()

		if !flagRefresh {
			return handleSegmentsLogic(ctx, idxClient)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleSegmentsLogic(ctx, idxClient)
		})
	},
}

func init() {
	getSegmentsCmd.Flags().StringVar(&flagSegmentsIndices, "indices", "", "Comma-separated list of indices")
}

func handleSegmentsLogic(ctx context.Context, idxClient index.Index) error {
	var indices []string
	if flagSegmentsIndices != "" {
		var err error
		indices, err = utils.ParseIndexPatternsCSV(flagSegmentsIndices, false)
		if err != nil {
			return err
		}
	}

	segments, err := idxClient.GetSegments(ctx, indices)
	if err != nil {
		return fmt.Errorf("failed to retrieve segments: %w", err)
	}

	return output.Render(segments)
}
