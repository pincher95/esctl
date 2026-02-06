package get

import (
	"github.com/pincher95/esctl/cmd/pipeline"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getPipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "List ingest pipelines",
	Example: utils.TrimAndIndent(`
	# List all pipelines
	esctl get pipelines

	# List pipelines by name substring
	esctl get pipelines --name logs
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleList(cmd.Context(), getPipelinesName)
	},
}

var getPipelinesName string

func init() {
	getPipelinesCmd.Flags().StringVar(&getPipelinesName, "name", "", "Filter pipelines by name or substring of pipeline name")
}
