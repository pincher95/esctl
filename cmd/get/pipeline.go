package get

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getPipelineCmd = &cobra.Command{
	Use:   "pipeline <pipeline-id>",
	Short: "Get an ingest pipeline",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get a specific pipeline
	esctl get pipeline my-pipeline
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleGet(cmd.Context(), args[0])
	},
}
