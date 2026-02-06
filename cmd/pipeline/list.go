package pipeline

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/pipeline"
	"github.com/pincher95/esctl/output"
)

func HandleList(ctx context.Context) error {
	pipelines, err := pipeline.ListPipelines(ctx)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}
	return output.Render(pipelines)
}
