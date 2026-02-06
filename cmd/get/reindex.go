package get

import (
	"github.com/pincher95/esctl/cmd/reindex"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getReindexCmd = &cobra.Command{
	Use:   "reindex <task-id>",
	Short: "Get reindex status",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get reindex status by task ID
	esctl get reindex task-id
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return reindex.HandleReindexStatus(cmd.Context(), args[0])
	},
}
