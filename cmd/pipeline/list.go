package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/pipeline"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
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
		pipelines = filtered
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(pipelines)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "NAME", Type: output.Text},
		{Header: "DESCRIPTION", Type: output.Text},
		{Header: "PROCESSORS", Type: output.Number},
		{Header: "VERSION", Type: output.Number},
	}

	data := make([][]string, 0, len(pipelines))
	for name, p := range pipelines {
		data = append(data, []string{
			name,
			p.Description,
			fmt.Sprintf("%d", len(p.Processors)),
			fmt.Sprintf("%d", p.Version),
		})
	}

	return output.PrintTable(columnDefs, data, nil)
}
