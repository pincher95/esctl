package get

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getSnapshotRepoCmd = &cobra.Command{
	Use:   "snapshot-repo <repository>",
	Short: "Get a snapshot repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get a snapshot repository
	esctl get snapshot-repo my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoGet(cmd.Context(), args[0])
	},
}
