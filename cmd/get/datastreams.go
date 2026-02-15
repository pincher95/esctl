package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/datastream"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
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

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(streams)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "NAME", Type: output.Text},
		{Header: "STATUS", Type: output.Text},
		{Header: "TEMPLATE", Type: output.Text},
		{Header: "ILM_POLICY", Type: output.Text},
		{Header: "GENERATION", Type: output.Number},
		{Header: "INDICES", Type: output.Number},
	}

	data := make([][]string, 0, len(streams))
	for _, s := range streams {
		data = append(data, []string{
			s.Name,
			s.Status,
			s.Template,
			s.ILMPolicy,
			fmt.Sprintf("%d", s.Generation),
			fmt.Sprintf("%d", len(s.Indices)),
		})
	}

	return output.PrintTable(columnDefs, data, nil)
}

func handleDataStreamLogic(ctx context.Context, name string) error {
	stream, err := datastream.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get data stream: %w", err)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(stream)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "setting", Type: output.Text},
		{Header: "value", Type: output.Text},
	}

	data := [][]string{
		{"name", stream.Name},
		{"status", stream.Status},
		{"template", stream.Template},
		{"timestamp_field", stream.TimestampField.Name},
		{"generation", fmt.Sprintf("%d", stream.Generation)},
		{"hidden", fmt.Sprintf("%t", stream.Hidden)},
	}
	if stream.ILMPolicy != "" {
		data = append(data, []string{"ilm_policy", stream.ILMPolicy})
	}
	for _, idx := range stream.Indices {
		data = append(data, []string{"backing_index", idx.IndexName})
	}

	return output.PrintTable(columnDefs, data, nil)
}
