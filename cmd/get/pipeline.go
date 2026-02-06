package get

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Get an ingest pipeline",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get a specific pipeline
	esctl get pipeline --id my-pipeline
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleGet(cmd.Context(), getPipelineID)
	},
}

var getPipelineID string

func init() {
	getPipelineCmd.Flags().StringVar(&getPipelineID, "id", "", "Pipeline ID")
	getPipelineCmd.MarkFlagRequired("id")
}
