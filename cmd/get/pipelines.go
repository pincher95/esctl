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
	Use:   "pipelines [ID]",
	Short: "Get or list ingest pipelines",
	Args:  cobra.MaximumNArgs(1),
	Example: utils.TrimAndIndent(`
	# List all pipelines
	esctl get pipelines

	# Get a specific pipeline by ID (positional argument)
	esctl get pipelines my-pipeline

	# Get a specific pipeline by ID (flag)
	esctl get pipelines --id my-pipeline

	# List pipelines by name substring
	esctl get pipelines --name logs
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Determine if we're getting a specific pipeline or listing all
		var pipelineID string
		if len(args) > 0 {
			pipelineID = args[0]
		} else if flagPipelinesID != "" {
			pipelineID = flagPipelinesID
		}

		// If a specific pipeline ID is provided, get that pipeline
		if pipelineID != "" {
			return handleGetSpecificPipeline(ctx, pipelineID)
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
	getPipelinesCmd.Flags().StringVar(&flagPipelinesID, "id", "", "Pipeline ID (for getting specific pipeline)")
}

func handleGetSpecificPipeline(ctx context.Context, id string) error {
	p, err := espipeline.GetPipeline(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %w", err)
	}
	return output.Render(p)
}
