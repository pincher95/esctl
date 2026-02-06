package reindex

import (
	"context"
	"fmt"

	esreindex "github.com/pincher95/esctl/es/reindex"
	"github.com/pincher95/esctl/output"
)

func HandleReindexStatus(ctx context.Context, taskID string) error {
	status, err := esreindex.GetReindexTaskStatus(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get reindex status: %w", err)
	}

	return output.Render(status)
}

func HandleReindexCancel(ctx context.Context, taskID string) error {
	if err := esreindex.CancelReindexTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to cancel reindex: %w", err)
	}

	fmt.Printf("Reindex task %s cancelled successfully\n", taskID)
	return nil
}
