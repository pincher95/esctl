package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/pipeline"
	"github.com/pincher95/esctl/output"
)

func HandleList(ctx context.Context, nameFilter string) error {
	pipelines, err := pipeline.ListPipelines(ctx)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}
	if nameFilter != "" {
		filtered := make(pipeline.PipelineResponse)
		for name, details := range pipelines {
			if strings.Contains(name, nameFilter) {
				filtered[name] = details
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("no pipelines matched: %s", nameFilter)
		}
		return output.Render(filtered)
	}
	return output.Render(pipelines)
}
