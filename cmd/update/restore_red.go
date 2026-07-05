package update

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/spf13/cobra"
)

var (
	restoreRedRepo            string
	restoreRedSnapshot        string
	restoreRedPattern         string
	restoreRedBatchSize       int
	restoreRedClose           bool
	restoreRedRenameAliasPat  string
	restoreRedRenameAliasRepl string
	restoreRedWait            bool
	restoreRedCMTimeout       string
	restoreRedDryRun          bool
)

var updateRestoreRedCmd = &cobra.Command{
	Use:   "restore-red",
	Short: "Restore red (unassigned-primary) indices from a snapshot, in batches",
	Long: utils.Trim(`
Recover indices whose primaries are unassigned (health red) by restoring them from a snapshot.

The command intersects the indices contained in the snapshot with the red indices matching
--pattern, so only recoverable, currently-broken indices are touched. Matches are restored in
batches: each batch is closed first (a closed index can be restored over) and then restored with
wait_for_completion so batches run sequentially rather than overwhelming the cluster manager.

This is the snapshot-based counterpart to the shard-recovery flow (get shard-stores / update
reroute): use it when no in-cluster copy survives but the data exists in a snapshot.`),
	Example: utils.TrimAndIndent(`
# Preview what would be restored
esctl update restore-red --repository my-repo --snapshot snap-1 --pattern "logz-*" --dry-run

# Restore red indices, renaming write-aliases so they don't collide with live ones
esctl update restore-red --repository my-repo --snapshot snap-1 --pattern "logz-*" \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias"
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		snapResp, err := snapshots.GetSnapshot(ctx, restoreRedRepo, restoreRedSnapshot)
		if err != nil {
			return err
		}
		if len(snapResp.Snapshots) == 0 {
			return fmt.Errorf("snapshot %q not found in repository %q", restoreRedSnapshot, restoreRedRepo)
		}
		inSnapshot := make(map[string]bool, len(snapResp.Snapshots[0].Indices))
		for _, idx := range snapResp.Snapshots[0].Indices {
			inSnapshot[idx] = true
		}

		indices, err := cat.NewCat().CatIndices(ctx, "", restoreRedPattern, "")
		if err != nil {
			return fmt.Errorf("failed to list indices for pattern %q: %w", restoreRedPattern, err)
		}

		var toRestore []string
		for _, i := range indices {
			if strings.EqualFold(i.Health, "red") && inSnapshot[i.Index] {
				toRestore = append(toRestore, i.Index)
			}
		}
		if len(toRestore) == 0 {
			fmt.Println("Nothing to restore: no red indices matching the pattern are present in the snapshot.")
			return nil
		}
		sort.Sort(sort.Reverse(sort.StringSlice(toRestore)))

		batches := batchStrings(toRestore, restoreRedBatchSize)
		fmt.Printf("Found %d red index(es) to restore from snapshot %q in %d batch(es) of up to %d.\n",
			len(toRestore), restoreRedSnapshot, len(batches), restoreRedBatchSize)

		idxClient := index.NewIndex()
		for n, b := range batches {
			joined := strings.Join(b, ",")
			if restoreRedDryRun {
				fmt.Printf("[dry-run] batch %d/%d (%d indices): would %srestore %s\n",
					n+1, len(batches), len(b), closeVerb(restoreRedClose), joined)
				continue
			}

			fmt.Printf("batch %d/%d: %d index(es)\n", n+1, len(batches), len(b))
			if restoreRedClose {
				if _, err := idxClient.Close(ctx, b); err != nil {
					return fmt.Errorf("batch %d/%d: failed to close indices: %w", n+1, len(batches), err)
				}
			}
			req := snapshots.RestoreSnapshotRequest{
				Indices:                joined,
				IncludeAliases:         true,
				RenameAliasPattern:     restoreRedRenameAliasPat,
				RenameAliasReplacement: restoreRedRenameAliasRepl,
			}
			if err := snapshots.RestoreSnapshotWithTimeout(ctx, restoreRedRepo, restoreRedSnapshot, req, restoreRedWait, restoreRedCMTimeout); err != nil {
				return fmt.Errorf("batch %d/%d: failed to restore: %w", n+1, len(batches), err)
			}
		}

		if restoreRedDryRun {
			fmt.Printf("Dry-run complete. Would restore %d red index(es).\n", len(toRestore))
		} else {
			fmt.Printf("Done. Restored %d red index(es).\n", len(toRestore))
		}
		return nil
	},
}

func closeVerb(closeFirst bool) string {
	if closeFirst {
		return "close + "
	}
	return ""
}

// batchStrings splits items into consecutive groups of at most size elements.
func batchStrings(items []string, size int) [][]string {
	if size <= 0 {
		size = 50
	}
	batches := make([][]string, 0, (len(items)+size-1)/size)
	for i := 0; i < len(items); i += size {
		batches = append(batches, items[i:min(i+size, len(items))])
	}
	return batches
}

func init() {
	updateRestoreRedCmd.Flags().StringVar(&restoreRedRepo, "repository", "", "Snapshot repository (required)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedSnapshot, "snapshot", "", "Snapshot name (required)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedPattern, "pattern", "", "Index pattern to match red indices (required)")
	updateRestoreRedCmd.Flags().IntVar(&restoreRedBatchSize, "batch-size", 50, "Number of indices to restore per batch")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedClose, "close", true, "Close each index before restoring it")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedRenameAliasPat, "rename-alias-pattern", "", "Regex to match alias names to rename on restore")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedRenameAliasRepl, "rename-alias-replacement", "", "Replacement for renamed aliases (may reference $1, $2, ...)")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedWait, "wait", true, "Wait for each batch to finish (wait_for_completion)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedCMTimeout, "cluster-manager-timeout", "5m", "Cluster-manager (master) timeout per restore request")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedDryRun, "dry-run", false, "Show what would be restored without changing anything")
	_ = updateRestoreRedCmd.MarkFlagRequired("repository")
	_ = updateRestoreRedCmd.MarkFlagRequired("snapshot")
	_ = updateRestoreRedCmd.MarkFlagRequired("pattern")

	updateCmd.AddCommand(updateRestoreRedCmd)
}
