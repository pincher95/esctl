package get

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getSnapshotStatusCmd = &cobra.Command{
	Use:   "snapshot-status",
	Short: "Get snapshot status",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get status of all snapshots
	esctl get snapshot-status

	# Get status for a repository
	esctl get snapshot-status --repository my-repo

	# Get status for a snapshot
	esctl get snapshot-status --repository my-repo --name my-snapshot
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if getSnapshotStatusName != "" && getSnapshotStatusRepo == "" {
			return fmt.Errorf("repository is required when using --name")
		}
		return runWithWatch(ctx, func() error {
			return snapshot.HandleSnapshotStatus(ctx, getSnapshotStatusRepo, getSnapshotStatusName)
		})
	},
}

var (
	getSnapshotStatusRepo string
	getSnapshotStatusName string
)

func init() {
	getSnapshotStatusCmd.Flags().StringVar(&getSnapshotStatusRepo, "repository", "", "Snapshot repository name")
	getSnapshotStatusCmd.Flags().StringVar(&getSnapshotStatusName, "name", "", "Snapshot name")
}
