package set

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var setPipelineFile string

var setPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Create or update an ingest pipeline",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Create pipeline from file
	esctl set pipeline --id my-pipeline --file=pipeline.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandlePut(cmd.Context(), setPipelineID, setPipelineFile)
	},
}

func init() {
	setPipelineCmd.Flags().StringVar(&setPipelineID, "id", "", "Pipeline ID")
	setPipelineCmd.Flags().StringVar(&setPipelineFile, "file", "", "JSON file containing pipeline definition")
	setPipelineCmd.MarkFlagRequired("id")
	setPipelineCmd.MarkFlagRequired("file")
}

var setPipelineID string
