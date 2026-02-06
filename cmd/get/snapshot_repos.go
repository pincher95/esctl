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
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoList(cmd.Context())
	},
}
