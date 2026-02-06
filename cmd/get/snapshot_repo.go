package get

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getSnapshotRepoCmd = &cobra.Command{
	Use:   "snapshot-repo",
	Short: "Get a snapshot repository",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get a snapshot repository
	esctl get snapshot-repo --repository my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoGet(cmd.Context(), getSnapshotRepoName)
	},
}

var getSnapshotRepoName string

func init() {
	getSnapshotRepoCmd.Flags().StringVar(&getSnapshotRepoName, "repository", "", "Snapshot repository name")
	getSnapshotRepoCmd.MarkFlagRequired("repository")
}
