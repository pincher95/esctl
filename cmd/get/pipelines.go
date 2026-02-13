package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	espipeline "github.com/pincher95/esctl/es/pipeline"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getPipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Get or list ingest pipelines",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# List all pipelines
	esctl get pipelines

	# Get a specific pipeline by ID
	esctl get pipelines --id my-pipeline

	# List pipelines by name substring
	esctl get pipelines --name logs
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// If a specific pipeline ID is provided, get that pipeline
		if flagPipelinesID != "" {
			return handleGetSpecificPipeline(ctx, flagPipelinesID)
		}

		// Otherwise, list all pipelines
		return pipeline.HandleList(ctx, flagPipelinesName)
	},
}

var (
	flagPipelinesName string
	flagPipelinesID   string
)

func init() {
	getPipelinesCmd.Flags().StringVar(&flagPipelinesName, "name", "", "Filter pipelines by name or substring of pipeline name")
	getPipelinesCmd.Flags().StringVar(&flagPipelinesID, "id", "", "Get a specific pipeline by ID")
}

func handleGetSpecificPipeline(ctx context.Context, id string) error {
	p, err := espipeline.GetPipeline(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %w", err)
	}
	return output.Render(p)
}
