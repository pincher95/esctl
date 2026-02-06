package set

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var setPipelineFile string

var setPipelineCmd = &cobra.Command{
	Use:   "pipeline <pipeline-id> --file=pipeline.json",
	Short: "Create or update an ingest pipeline",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Create pipeline from file
	esctl set pipeline my-pipeline --file=pipeline.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandlePut(cmd.Context(), args[0], setPipelineFile)
	},
}

func init() {
	setPipelineCmd.Flags().StringVar(&setPipelineFile, "file", "", "JSON file containing pipeline definition")
	setPipelineCmd.MarkFlagRequired("file")
}
