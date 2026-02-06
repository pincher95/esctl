package delete

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Delete a snapshot",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Delete snapshot
	esctl delete snapshot --repository my-repo --name my-snapshot
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleSnapshotDelete(cmd.Context(), deleteSnapshotRepo, deleteSnapshotName)
	},
}

var (
	deleteSnapshotRepo string
	deleteSnapshotName string
)

func init() {
	deleteSnapshotCmd.Flags().StringVar(&deleteSnapshotRepo, "repository", "", "Snapshot repository name")
	deleteSnapshotCmd.Flags().StringVar(&deleteSnapshotName, "name", "", "Snapshot name")
	deleteSnapshotCmd.MarkFlagRequired("repository")
	deleteSnapshotCmd.MarkFlagRequired("name")
}
