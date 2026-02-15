package set

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/utils"
	esreindex "github.com/pincher95/esctl/es/reindex"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	setReindexSource            string
	setReindexDest              string
	setReindexQuery             string
	setReindexScript            string
	setReindexWait              bool
	setReindexRequestsPerSecond float64
	setReindexTimeout           string
	setReindexRefresh           bool
	setReindexSlices            string
	setReindexConflicts         string
	setReindexSize              int
)

var setReindexCmd = &cobra.Command{
	Use:   "reindex --source=<index> --dest=<index>",
	Short: "Start a reindex operation",
	Example: utils.TrimAndIndent(`
	# Basic reindex from one index to another
	esctl set reindex --source=old-index --dest=new-index

	# Reindex with query filter
	esctl set reindex --source=logs --dest=filtered-logs --query='{"range":{"@timestamp":{"gte":"2023-01-01"}}}'

	# Reindex with script transformation
	esctl set reindex --source=users --dest=users-v2 --script='ctx._source.full_name = ctx._source.first_name + " " + ctx._source.last_name'

	# Reindex with rate limiting and wait for completion
	esctl set reindex --source=big-index --dest=new-index --requests-per-second=100 --wait
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleSetReindexStart(cmd.Context())
	},
}

func init() {
	setReindexCmd.Flags().StringVar(&setReindexSource, "source", "", "Source index name")
	setReindexCmd.Flags().StringVar(&setReindexDest, "dest", "", "Destination index name")
	setReindexCmd.Flags().StringVar(&setReindexQuery, "query", "", "Query to filter documents (JSON)")
	setReindexCmd.Flags().StringVar(&setReindexScript, "script", "", "Painless script for document transformation")
	setReindexCmd.Flags().BoolVar(&setReindexWait, "wait", false, "Wait for reindex completion")
	setReindexCmd.Flags().Float64Var(&setReindexRequestsPerSecond, "requests-per-second", -1, "Throttle rate in requests per second")
	setReindexCmd.Flags().StringVar(&setReindexTimeout, "timeout", "", "Timeout for the operation")
	setReindexCmd.Flags().BoolVar(&setReindexRefresh, "refresh", false, "Refresh the destination index")
	setReindexCmd.Flags().StringVar(&setReindexSlices, "slices", "1", "Number of slices for parallel processing")
	setReindexCmd.Flags().StringVar(&setReindexConflicts, "conflicts", "abort", "How to handle version conflicts (abort, proceed)")
	setReindexCmd.Flags().IntVar(&setReindexSize, "size", 0, "Number of documents to process per batch")

	_ = setReindexCmd.MarkFlagRequired("source")
	_ = setReindexCmd.MarkFlagRequired("dest")
}

func handleSetReindexStart(ctx context.Context) error {
	if err := validation.ValidateIndexPattern(setReindexSource); err != nil {
		return err
	}
	if err := validation.ValidateIndexName(setReindexDest); err != nil {
		return err
	}
	if setReindexRequestsPerSecond < -1 {
		return fmt.Errorf("requests-per-second must be -1 or >= 0, got %f", setReindexRequestsPerSecond)
	}
	if setReindexTimeout != "" {
		if err := validation.ValidateTimeout(setReindexTimeout); err != nil {
			return err
		}
	}
	if setReindexConflicts != "abort" && setReindexConflicts != "proceed" {
		return fmt.Errorf("invalid conflicts value: %s (must be 'abort' or 'proceed')", setReindexConflicts)
	}
	if setReindexSize < 0 {
		return fmt.Errorf("size cannot be negative: %d", setReindexSize)
	}

	source := esreindex.ReindexSource{
		Index: setReindexSource,
	}

	if setReindexQuery != "" {
		queryMap, err := utils.ParseJSONMap(setReindexQuery, "invalid query JSON")
		if err != nil {
			return err
		}
		source.Query = queryMap
	}

	if setReindexSize > 0 {
		source.Size = &setReindexSize
	}

	dest := esreindex.ReindexDest{
		Index: setReindexDest,
	}

	request := esreindex.ReindexRequest{
		Source:    source,
		Dest:      dest,
		Conflicts: setReindexConflicts,
	}

	if setReindexScript != "" {
		request.Script = &esreindex.ReindexScript{
			Source: setReindexScript,
			Lang:   "painless",
		}
	}

	var slices any
	if setReindexSlices != "" {
		if slicesInt, err := strconv.Atoi(setReindexSlices); err == nil {
			if slicesInt <= 0 {
				return fmt.Errorf("slices must be greater than 0, got %d", slicesInt)
			}
			slices = slicesInt
		} else {
			slices = setReindexSlices
		}
	}

	result, err := esreindex.StartReindex(ctx, request, setReindexWait, setReindexRequestsPerSecond, setReindexTimeout, setReindexRefresh, slices)
	if err != nil {
		return fmt.Errorf("failed to start reindex: %w", err)
	}

	if result.Task != "" {
		fmt.Printf("Reindex task started: %s\n", result.Task)
		fmt.Printf("Use 'esctl get reindex %s' to check progress\n", result.Task)
	} else {
		fmt.Printf("Reindex completed: %d documents processed\n", result.Total)
		if result.VersionConflicts > 0 {
			fmt.Printf("Version conflicts: %d\n", result.VersionConflicts)
		}
	}

	return output.Render(result)
}
