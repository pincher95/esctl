package get

import (
	"github.com/pincher95/esctl/cmd/reindex"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Get reindex status",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get reindex status by task ID
	esctl get reindex --task-id task-id
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return reindex.HandleReindexStatus(cmd.Context(), getReindexTaskID)
	},
}

var getReindexTaskID string

func init() {
	getReindexCmd.Flags().StringVar(&getReindexTaskID, "task-id", "", "Reindex task ID")
	getReindexCmd.MarkFlagRequired("task-id")
}
