package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/datastream"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flagDataStreamsName string

var getDataStreamsCmd = &cobra.Command{
	Use:   "data-streams",
	Short: "List data streams or get a specific data stream",
	Long: utils.Trim(`
		Lists all data streams in the Elasticsearch cluster, or gets a specific data stream by name.
		Data streams are append-only time series indices that automatically manage backing indices.
	`),
	Example: utils.TrimAndIndent(`
		# List all data streams
		esctl get data-streams

		# Get a specific data stream by name
		esctl get data-streams --name logs-app

		# List data streams matching a pattern
		esctl get data-streams --name "logs-*"

		# Watch mode for monitoring
		esctl get data-streams --watch

		# JSON output
		esctl get data-streams -o json
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// If specific name provided (without wildcards), get specific data stream
		// Otherwise list data streams (with optional pattern filter)
		if flagDataStreamsName != "" && !strings.ContainsAny(flagDataStreamsName, "*?") {
			if !flagRefresh {
				return handleDataStreamLogic(ctx, flagDataStreamsName)
			}
			return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
				return handleDataStreamLogic(ctx, flagDataStreamsName)
			})
		}

		// List data streams (with optional pattern)
		if !flagRefresh {
			return handleDataStreamsLogic(ctx, flagDataStreamsName)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleDataStreamsLogic(ctx, flagDataStreamsName)
		})
	},
}

func init() {
	getDataStreamsCmd.Flags().StringVar(&flagDataStreamsName, "name", "", "Data stream name or pattern to retrieve")
}

func handleDataStreamsLogic(ctx context.Context, name string) error {
	streams, err := datastream.List(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to list data streams: %w", err)
	}

	return output.Render(streams)
}

func handleDataStreamLogic(ctx context.Context, name string) error {
	stream, err := datastream.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	return output.Render(stream)
}
