package pipeline

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/pipeline"
	"github.com/pincher95/esctl/output"
)

func HandleGet(ctx context.Context, id string) error {
	p, err := pipeline.GetPipeline(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %w", err)
	}
	return output.Render(p)
}
