package delete

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteSnapshotRepoCmd = &cobra.Command{
	Use:   "snapshot-repo",
	Short: "Delete a snapshot repository",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Delete a repository
	esctl delete snapshot-repo --repository my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoDelete(cmd.Context(), deleteSnapshotRepoName)
	},
}

var deleteSnapshotRepoName string

func init() {
	deleteSnapshotRepoCmd.Flags().StringVar(&deleteSnapshotRepoName, "repository", "", "Snapshot repository name")
	deleteSnapshotRepoCmd.MarkFlagRequired("repository")
}
