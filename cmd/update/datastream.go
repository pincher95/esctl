package update

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/datastream"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flagDataStreamName string

var updateDataStreamCmd = &cobra.Command{
	Use:   "data-stream",
	Short: "Perform operations on data streams",
	Long: utils.Trim(`
		Performs update operations on data streams, such as rolling over to a new backing index.
	`),
}

var rolloverCmd = &cobra.Command{
	Use:   "rollover",
	Short: "Rollover a data stream to a new backing index",
	Long: utils.Trim(`
		Creates a new backing index for the data stream and makes it the write index.
		The old write index becomes a read-only backing index.
	`),
	Example: utils.TrimAndIndent(`
		# Rollover a data stream
		esctl update data-stream rollover --name logs-app

		# View rollover result in JSON
		esctl update data-stream rollover --name logs-app -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleRollover(ctx, flagDataStreamName)
	},
}

func init() {
	rolloverCmd.Flags().StringVar(&flagDataStreamName, "name", "", "Data stream name")
	_ = rolloverCmd.MarkFlagRequired("name")
	updateDataStreamCmd.AddCommand(rolloverCmd)
}

func handleRollover(ctx context.Context, name string) error {
	resp, err := datastream.Rollover(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to rollover data stream: %w", err)
	}

	if resp.RolledOver {
		fmt.Printf("Successfully rolled over data stream '%s'\n", name)
		fmt.Printf("  Old index: %s\n", resp.OldIndex)
		fmt.Printf("  New index: %s\n", resp.NewIndex)
	} else {
		fmt.Printf("Data stream '%s' did not rollover\n", name)
	}

	return output.Render(resp)
}
