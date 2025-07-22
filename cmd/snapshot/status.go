package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	statusRepository string
	statusSnapshot   string
)

var statusCmd = &cobra.Command{
	Use:   "status [repository] [snapshot]",
	Short: "Get snapshot status",
	Args:  cobra.RangeArgs(0, 2),
	Example: utils.TrimAndIndent(`
	# Get status of all snapshots
	esctl snapshot status

	# Get status of all snapshots in a repository
	esctl snapshot status my-repo

	# Get status of a specific snapshot
	esctl snapshot status my-repo my-snapshot
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

		return handleSnapshotStatus(ctx, repo, snap)
	},
}

func handleSnapshotStatus(ctx context.Context, repository, snapshot string) error {
	result, err := snapshots.SnapshotStatus(ctx, repository, snapshot)
	if err != nil {
		return fmt.Errorf("failed to get snapshot status: %w", err)
	}

	output.PrintJson(result)
	return nil
}
