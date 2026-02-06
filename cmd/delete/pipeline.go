package delete

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deletePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Delete an ingest pipeline",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Delete pipeline
	esctl delete pipeline --id my-pipeline
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleDelete(cmd.Context(), deletePipelineID)
	},
}

var deletePipelineID string

func init() {
	deletePipelineCmd.Flags().StringVar(&deletePipelineID, "id", "", "Pipeline ID")
	deletePipelineCmd.MarkFlagRequired("id")
}
