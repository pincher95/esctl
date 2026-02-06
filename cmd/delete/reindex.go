package delete

import (
	"github.com/pincher95/esctl/cmd/reindex"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteReindexCmd = &cobra.Command{
	Use:   "reindex <task-id>",
	Short: "Cancel a reindex task",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Cancel reindex task
	esctl delete reindex task-id
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return reindex.HandleReindexCancel(cmd.Context(), args[0])
	},
}
