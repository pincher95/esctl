package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/spf13/cobra"
)

var (
	restoreIndices             string
	restoreWait                bool
	restoreIgnoreUnavailable   bool
	restoreIncludeGlobalState  bool
	restoreRenamePattern       string
	restoreRenameReplacement   string
	restoreIncludeAliases      bool
	restoreIndexSettings       string
	restoreIgnoreIndexSettings string
)

var restoreCmd = &cobra.Command{
	Use:   "restore <repository> <snapshot>",
	Short: "Restore a snapshot",
	Args:  cobra.ExactArgs(2),
	Example: utils.TrimAndIndent(`
	# Restore all indices from a snapshot
	esctl snapshot restore my-repo my-snapshot

	# Restore specific indices
	esctl snapshot restore my-repo my-snapshot --indices="index1,index2"

	# Restore with rename pattern
	esctl snapshot restore my-repo my-snapshot --rename-pattern="(.+)" --rename-replacement="restored_$1"

	# Restore and wait for completion
	esctl snapshot restore my-repo my-snapshot --wait
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSnapshotRestore(ctx, args[0], args[1])
	},
}

func init() {
	restoreCmd.Flags().StringVar(&restoreIndices, "indices", "", "Comma-separated list of indices to restore")
	restoreCmd.Flags().BoolVar(&restoreWait, "wait", false, "Wait for restore completion")
	restoreCmd.Flags().BoolVar(&restoreIgnoreUnavailable, "ignore-unavailable", false, "Ignore unavailable indices")
	restoreCmd.Flags().BoolVar(&restoreIncludeGlobalState, "include-global-state", false, "Include global cluster state")
	restoreCmd.Flags().StringVar(&restoreRenamePattern, "rename-pattern", "", "Rename pattern for restored indices")
	restoreCmd.Flags().StringVar(&restoreRenameReplacement, "rename-replacement", "", "Rename replacement for restored indices")
	restoreCmd.Flags().BoolVar(&restoreIncludeAliases, "include-aliases", true, "Include aliases when restoring")
	restoreCmd.Flags().StringVar(&restoreIndexSettings, "index-settings", "", "Index settings to apply during restore (key:value pairs)")
	restoreCmd.Flags().StringVar(&restoreIgnoreIndexSettings, "ignore-index-settings", "", "Comma-separated list of index settings to ignore")
}

func handleSnapshotRestore(ctx context.Context, repository, snapshot string) error {
	request := snapshots.RestoreSnapshotRequest{
		Indices:            restoreIndices,
		IgnoreUnavailable:  restoreIgnoreUnavailable,
		IncludeGlobalState: restoreIncludeGlobalState,
		RenamePattern:      restoreRenamePattern,
		RenameReplacement:  restoreRenameReplacement,
		IncludeAliases:     restoreIncludeAliases,
	}

	// Parse index settings if provided
	if restoreIndexSettings != "" {
		settings, err := parseSettings(restoreIndexSettings)
		if err != nil {
			return fmt.Errorf("invalid index settings format: %w", err)
		}
		request.IndexSettings = settings
	}

	// Parse ignored index settings
	if restoreIgnoreIndexSettings != "" {
		request.IgnoreIndexSettings = strings.Split(restoreIgnoreIndexSettings, ",")
		for i, setting := range request.IgnoreIndexSettings {
			request.IgnoreIndexSettings[i] = strings.TrimSpace(setting)
		}
	}

	err := snapshots.RestoreSnapshot(ctx, repository, snapshot, request, restoreWait)
	if err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	if restoreWait {
		fmt.Printf("Snapshot '%s' restored successfully from repository '%s'\n", snapshot, repository)
	} else {
		fmt.Printf("Snapshot '%s' restore started from repository '%s'\n", snapshot, repository)
	}

	return nil
}
