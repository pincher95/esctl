package reindex

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/reindex"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex documents between indices",
	Long:  "Start, monitor, and control reindex operations in Elasticsearch",
}

var (
	startSource            string
	startDest              string
	startQuery             string
	startScript            string
	startWait              bool
	startRequestsPerSecond float64
	startTimeout           string
	startRefresh           bool
	startSlices            string
	startConflicts         string
	startSize              int
)

var startCmd = &cobra.Command{
	Use:   "start --source=<index> --dest=<index>",
	Short: "Start a reindex operation",
	Example: utils.TrimAndIndent(`
	# Basic reindex from one index to another
	esctl reindex start --source=old-index --dest=new-index

	# Reindex with query filter
	esctl reindex start --source=logs --dest=filtered-logs --query='{"range":{"@timestamp":{"gte":"2023-01-01"}}}'

	# Reindex with script transformation
	esctl reindex start --source=users --dest=users-v2 --script='ctx._source.full_name = ctx._source.first_name + " " + ctx._source.last_name'

	# Reindex with rate limiting and wait for completion
	esctl reindex start --source=big-index --dest=new-index --requests-per-second=100 --wait
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleReindexStart(ctx)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status <task-id>",
	Short: "Get status of a reindex task",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get status of a reindex task
	esctl reindex status r1A2WoRbTwKZ516z6NEs5A:36619
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleReindexStatus(ctx, args[0])
	},
}

var cancelCmd = &cobra.Command{
	Use:   "cancel <task-id>",
	Short: "Cancel a reindex task",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Cancel a reindex task
	esctl reindex cancel r1A2WoRbTwKZ516z6NEs5A:36619
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleReindexCancel(ctx, args[0])
	},
}

func init() {
	startCmd.Flags().StringVar(&startSource, "source", "", "Source index name")
	startCmd.Flags().StringVar(&startDest, "dest", "", "Destination index name")
	startCmd.Flags().StringVar(&startQuery, "query", "", "Query to filter documents (JSON)")
	startCmd.Flags().StringVar(&startScript, "script", "", "Painless script for document transformation")
	startCmd.Flags().BoolVar(&startWait, "wait", false, "Wait for reindex completion")
	startCmd.Flags().Float64Var(&startRequestsPerSecond, "requests-per-second", -1, "Throttle rate in requests per second")
	startCmd.Flags().StringVar(&startTimeout, "timeout", "", "Timeout for the operation")
	startCmd.Flags().BoolVar(&startRefresh, "refresh", false, "Refresh the destination index")
	startCmd.Flags().StringVar(&startSlices, "slices", "1", "Number of slices for parallel processing")
	startCmd.Flags().StringVar(&startConflicts, "conflicts", "abort", "How to handle version conflicts (abort, proceed)")
	startCmd.Flags().IntVar(&startSize, "size", 0, "Number of documents to process per batch")

	startCmd.MarkFlagRequired("source")
	startCmd.MarkFlagRequired("dest")

	reindexCmd.AddCommand(startCmd)
	reindexCmd.AddCommand(statusCmd)
	reindexCmd.AddCommand(cancelCmd)
}

func Cmd() *cobra.Command {
	return reindexCmd
}

func handleReindexStart(ctx context.Context) error {
	// Build source
	source := reindex.ReindexSource{
		Index: startSource,
	}

	// Parse query if provided
	if startQuery != "" {
		var queryMap map[string]interface{}
		if err := json.Unmarshal([]byte(startQuery), &queryMap); err != nil {
			return fmt.Errorf("invalid query JSON: %w", err)
		}
		source.Query = queryMap
	}

	if startSize > 0 {
		source.Size = &startSize
	}

	// Build destination
	dest := reindex.ReindexDest{
		Index: startDest,
	}

	// Build request
	request := reindex.ReindexRequest{
		Source:    source,
		Dest:      dest,
		Conflicts: startConflicts,
	}

	// Add script if provided
	if startScript != "" {
		request.Script = &reindex.ReindexScript{
			Source: startScript,
			Lang:   "painless",
		}
	}

	// Parse slices
	var slices interface{}
	if startSlices != "" {
		if slicesInt, err := strconv.Atoi(startSlices); err == nil {
			slices = slicesInt
		} else {
			slices = startSlices
		}
	}

	result, err := reindex.StartReindex(ctx, request, startWait, startRequestsPerSecond, startTimeout, startRefresh, slices)
	if err != nil {
		return fmt.Errorf("failed to start reindex: %w", err)
	}

	if result.Task != "" {
		fmt.Printf("Reindex task started: %s\n", result.Task)
		fmt.Printf("Use 'esctl reindex status %s' to check progress\n", result.Task)
	} else {
		fmt.Printf("Reindex completed: %d documents processed\n", result.Total)
		if result.VersionConflicts > 0 {
			fmt.Printf("Version conflicts: %d\n", result.VersionConflicts)
		}
	}

	output.PrintJson(result)
	return nil
}

func handleReindexStatus(ctx context.Context, taskID string) error {
	status, err := reindex.GetReindexTaskStatus(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get reindex status: %w", err)
	}

	output.PrintJson(status)
	return nil
}

func handleReindexCancel(ctx context.Context, taskID string) error {
	err := reindex.CancelReindexTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to cancel reindex: %w", err)
	}

	fmt.Printf("Reindex task %s cancelled successfully\n", taskID)
	return nil
}
