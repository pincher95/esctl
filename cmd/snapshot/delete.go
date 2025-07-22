package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <repository> <snapshot>",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(2),
	Example: utils.TrimAndIndent(`
	# Delete a snapshot
	esctl snapshot delete my-repo my-snapshot
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSnapshotDelete(ctx, args[0], args[1])
	},
}

func handleSnapshotDelete(ctx context.Context, repository, snapshot string) error {
	err := snapshots.DeleteSnapshot(ctx, repository, snapshot)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	fmt.Printf("Snapshot '%s' deleted successfully from repository '%s'\n", snapshot, repository)
	return nil
}
