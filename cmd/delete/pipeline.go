package delete

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deletePipelineCmd = &cobra.Command{
	Use:   "pipeline <pipeline-id>",
	Short: "Delete an ingest pipeline",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Delete pipeline
	esctl delete pipeline my-pipeline
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleDelete(cmd.Context(), args[0])
	},
}
