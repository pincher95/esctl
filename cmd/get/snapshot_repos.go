package get

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getSnapshotReposCmd = &cobra.Command{
	Use:   "snapshot-repos",
	Short: "List snapshot repositories",
	Example: utils.TrimAndIndent(`
	# List snapshot repositories
	esctl get snapshot-repos

	# List snapshot repositories by name substring
	esctl get snapshot-repos --name backup
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoList(cmd.Context(), getSnapshotReposName)
	},
}

var getSnapshotReposName string

func init() {
	getSnapshotReposCmd.Flags().StringVar(&getSnapshotReposName, "name", "", "Filter repositories by name or substring of repository name")
}
