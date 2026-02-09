package pipeline

import (
	"context"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/pipeline"
	"github.com/pincher95/esctl/output"
)

func HandleSimulate(ctx context.Context, filePath, pipelineID string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var req pipeline.SimulateRequest
	if err := utils.UnmarshalJSON(data, &req, "invalid simulation request JSON"); err != nil {
		return err
	}

	result, err := pipeline.SimulatePipeline(ctx, pipelineID, req)
	if err != nil {
		return fmt.Errorf("failed to simulate pipeline: %w", err)
	}

	return output.Render(result)
}
