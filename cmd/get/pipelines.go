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
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pipeline.HandleList(cmd.Context())
	},
}
