package pipeline

import (
	"context"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/pipeline"
)

func HandlePut(ctx context.Context, id, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var p pipeline.Pipeline
	if err := utils.UnmarshalJSON(data, &p, "invalid pipeline JSON"); err != nil {
		return err
	}

	if err := pipeline.PutPipeline(ctx, id, p); err != nil {
		return fmt.Errorf("failed to put pipeline: %w", err)
	}

	fmt.Printf("Pipeline '%s' created/updated successfully\n", id)
	return nil
}
