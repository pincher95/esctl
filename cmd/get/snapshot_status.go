package get

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getSnapshotStatusCmd = &cobra.Command{
	Use:   "snapshot-status [repository] [snapshot]",
	Short: "Get snapshot status",
	Args:  cobra.RangeArgs(0, 2),
	Example: utils.TrimAndIndent(`
	# Get status of all snapshots
	esctl get snapshot-status

	# Get status for a repository
	esctl get snapshot-status my-repo

	# Get status for a snapshot
	esctl get snapshot-status my-repo my-snapshot
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var repo, snap string
		if len(args) > 0 {
			repo = args[0]
		}
		if len(args) > 1 {
			snap = args[1]
		}
		return snapshot.HandleSnapshotStatus(ctx, repo, snap)
	},
}
