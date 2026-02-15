package update

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	updatePipelineFile string
	updatePipelineID   string
)

var updatePipelineCmd = &cobra.Command{
	Use:   "pipeline --file=request.json [--pipeline=<pipeline-id>]",
	Short: "Simulate pipeline execution",
	Example: utils.TrimAndIndent(`
	# Simulate existing pipeline
	esctl update pipeline --pipeline=my-pipeline --file=docs.json

	# Simulate inline pipeline definition
	esctl update pipeline --file=pipeline-with-docs.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleSimulate(cmd.Context(), updatePipelineFile, updatePipelineID)
	},
}

func init() {
	updatePipelineCmd.Flags().StringVar(&updatePipelineFile, "file", "", "JSON file containing simulation request")
	updatePipelineCmd.Flags().StringVar(&updatePipelineID, "pipeline", "", "Pipeline ID to simulate (optional)")
	updatePipelineCmd.MarkFlagRequired("file")
}
