package pipeline

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/pipeline"
)

func HandleDelete(ctx context.Context, id string) error {
	if err := pipeline.DeletePipeline(ctx, id); err != nil {
		return fmt.Errorf("failed to delete pipeline: %w", err)
	}
	fmt.Printf("Pipeline '%s' deleted successfully\n", id)
	return nil
}
