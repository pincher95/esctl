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
	flagRecoveryIndices  string
	flagRecoveryDetailed bool
)

var getRecoveryCmd = &cobra.Command{
	Use:   "recovery",
	Short: "Get shard recovery information",
	Long: utils.Trim(`
		Retrieves information about ongoing and completed shard recoveries. This is useful for
		monitoring shard allocation, rebalancing, and understanding recovery performance.
	`),
	Example: utils.TrimAndIndent(`
		# Get recovery info for all indices
		esctl get recovery

		# Get recovery info for specific indices
		esctl get recovery --indices my-index,other-index

		# Get detailed recovery information
		esctl get recovery --indices my-index --detailed

		# Watch mode for monitoring active recoveries
		esctl get recovery --watch --interval 5s

		# JSON output
		esctl get recovery --indices my-index -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idxClient := index.NewIndex()

		if !flagRefresh {
			return handleRecoveryLogic(ctx, idxClient)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleRecoveryLogic(ctx, idxClient)
		})
	},
}

func init() {
	getRecoveryCmd.Flags().StringVar(&flagRecoveryIndices, "indices", "", "Comma-separated list of indices")
	getRecoveryCmd.Flags().BoolVar(&flagRecoveryDetailed, "detailed", false, "Include detailed recovery information")
}

func handleRecoveryLogic(ctx context.Context, idxClient index.Index) error {
	var indices []string
	if flagRecoveryIndices != "" {
		var err error
		indices, err = utils.ParseIndexPatternsCSV(flagRecoveryIndices, false)
		if err != nil {
			return err
		}
	}

	recovery, err := idxClient.GetRecovery(ctx, indices, flagRecoveryDetailed)
	if err != nil {
		return fmt.Errorf("failed to retrieve recovery info: %w", err)
	}

	return output.Render(recovery)
}
