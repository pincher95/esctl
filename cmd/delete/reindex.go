package delete

import (
	"github.com/pincher95/esctl/cmd/reindex"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Cancel a reindex task",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Cancel reindex task
	esctl delete reindex --task-id task-id
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return reindex.HandleReindexCancel(cmd.Context(), deleteReindexTaskID)
	},
}

var deleteReindexTaskID string

func init() {
	deleteReindexCmd.Flags().StringVar(&deleteReindexTaskID, "task-id", "", "Reindex task ID")
	deleteReindexCmd.MarkFlagRequired("task-id")
}
