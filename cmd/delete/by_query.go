package delete

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/delete"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	byQueryIndices           string
	byQueryQuery             string
	byQueryMaxDocs           int
	byQueryConflicts         string
	byQueryRefresh           bool
	byQueryTimeout           string
	byQueryWait              bool
	byQueryRequestsPerSecond float64
	byQuerySlices            string
)

var byQueryCmd = &cobra.Command{
	Use:   "by-query",
	Short: "Delete documents by query",
	Example: utils.TrimAndIndent(`
	# Delete all documents in an index
	esctl delete by-query --indices="my-index" --query='{"match_all":{}}'

	# Delete documents matching a specific query
	esctl delete by-query --indices="logs-*" --query='{"range":{"@timestamp":{"lt":"2023-01-01"}}}'

	# Delete with conflicts handling
	esctl delete by-query --indices="my-index" --query='{"term":{"status":"old"}}' --conflicts=proceed

	# Delete with rate limiting
	esctl delete by-query --indices="my-index" --query='{"match_all":{}}' --requests-per-second=100
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleByQueryDelete(ctx)
	},
}

func init() {
	byQueryCmd.Flags().StringVar(&byQueryIndices, "indices", "", "Comma-separated list of indices to search")
	byQueryCmd.Flags().StringVar(&byQueryQuery, "query", "", "Query to match documents for deletion (JSON)")
	byQueryCmd.Flags().IntVar(&byQueryMaxDocs, "max-docs", 0, "Maximum number of documents to delete")
	byQueryCmd.Flags().StringVar(&byQueryConflicts, "conflicts", "abort", "How to handle version conflicts (abort, proceed)")
	byQueryCmd.Flags().BoolVar(&byQueryRefresh, "refresh", false, "Refresh the affected shards")
	byQueryCmd.Flags().StringVar(&byQueryTimeout, "timeout", "", "Timeout for the request")
	byQueryCmd.Flags().BoolVar(&byQueryWait, "wait", false, "Wait for the operation to complete")
	byQueryCmd.Flags().Float64Var(&byQueryRequestsPerSecond, "requests-per-second", -1, "Throttle rate in requests per second")
	byQueryCmd.Flags().StringVar(&byQuerySlices, "slices", "1", "Number of slices for parallel processing")

	byQueryCmd.MarkFlagRequired("query")
}

func handleByQueryDelete(ctx context.Context) error {
	if byQueryConflicts != "abort" && byQueryConflicts != "proceed" {
		return fmt.Errorf("invalid conflicts value: %s (must be 'abort' or 'proceed')", byQueryConflicts)
	}
	if byQueryMaxDocs < 0 {
		return fmt.Errorf("max-docs cannot be negative: %d", byQueryMaxDocs)
	}
	if byQueryRequestsPerSecond < -1 {
		return fmt.Errorf("requests-per-second must be -1 or >= 0, got %f", byQueryRequestsPerSecond)
	}
	if byQueryTimeout != "" {
		if err := validation.ValidateTimeout(byQueryTimeout); err != nil {
			return err
		}
	}

	// Parse query JSON
	queryMap, err := utils.ParseJSONMap(byQueryQuery, "invalid query JSON")
	if err != nil {
		return err
	}

	// Parse indices
	var indices []string
	if byQueryIndices != "" {
		parsed, err := utils.ParseIndexPatternsCSV(byQueryIndices, true)
		if err != nil {
			return err
		}
		indices = parsed
	}

	// Build request
	request := delete.DeleteByQueryRequest{
		Query:             queryMap,
		Conflicts:         byQueryConflicts,
		Refresh:           byQueryRefresh,
		Timeout:           byQueryTimeout,
		WaitForCompletion: byQueryWait,
	}

	if byQueryMaxDocs > 0 {
		request.MaxDocs = &byQueryMaxDocs
	}

	if byQueryRequestsPerSecond >= 0 {
		request.RequestsPerSecond = &byQueryRequestsPerSecond
	}

	// Parse slices
	if byQuerySlices != "" {
		if slicesInt, err := strconv.Atoi(byQuerySlices); err == nil {
			if slicesInt <= 0 {
				return fmt.Errorf("slices must be greater than 0, got %d", slicesInt)
			}
			request.Slices = slicesInt
		} else {
			request.Slices = byQuerySlices
		}
	}

	result, err := delete.DeleteByQuery(ctx, indices, request)
	if err != nil {
		return fmt.Errorf("failed to delete by query: %w", err)
	}

	if result.Task != "" {
		fmt.Printf("Delete by query task started: %s\n", result.Task)
	} else {
		fmt.Printf("Deleted %d documents out of %d total\n", result.Deleted, result.Total)
		if result.VersionConflicts > 0 {
			fmt.Printf("Version conflicts: %d\n", result.VersionConflicts)
		}
	}

	return output.Render(result)
}
