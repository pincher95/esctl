package delete

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteSnapshotCmd = &cobra.Command{
	Use:   "snapshot <repository> <snapshot>",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(2),
	Example: utils.TrimAndIndent(`
	# Delete snapshot
	esctl delete snapshot my-repo my-snapshot
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleSnapshotDelete(cmd.Context(), args[0], args[1])
	},
}
