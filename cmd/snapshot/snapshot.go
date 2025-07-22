package snapshot

import (
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage Elasticsearch snapshots and repositories",
	Long:  "Create, list, delete snapshots and manage snapshot repositories in Elasticsearch",
}

func Cmd() *cobra.Command {
	snapshotCmd.AddCommand(repoCmd)
	snapshotCmd.AddCommand(createCmd)
	snapshotCmd.AddCommand(statusCmd)
	snapshotCmd.AddCommand(restoreCmd)
	snapshotCmd.AddCommand(listCmd)
	snapshotCmd.AddCommand(deleteCmd)
	return snapshotCmd
}
