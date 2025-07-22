package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <repository>",
	Short: "List snapshots in a repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# List all snapshots in a repository
	esctl snapshot list my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSnapshotList(ctx, args[0])
	},
}

func handleSnapshotList(ctx context.Context, repository string) error {
	result, err := snapshots.ListSnapshots(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}

	output.PrintJson(result)
	return nil
}
