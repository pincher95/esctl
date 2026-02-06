package delete

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteSnapshotRepoCmd = &cobra.Command{
	Use:   "snapshot-repo <repository>",
	Short: "Delete a snapshot repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Delete a repository
	esctl delete snapshot-repo my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoDelete(cmd.Context(), args[0])
	},
}
