package update

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

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
	restoreRedAliasPattern    string
	restoreRedExcludeToday    bool
	restoreRedReplicas        int
	restoreRedBoxType         string
	restoreRedIgnoreSettings  []string
	restoreRedBatchSize       int
	restoreRedClose           bool
	restoreRedRenameAliasPat  string
	restoreRedRenameAliasRepl string
	restoreRedWait            bool
	restoreRedCMTimeout       string
	restoreRedDryRun          bool
	restoreRedIncludeClosed   bool
	restoreRedIncludeAliases  bool
	restoreRedPollInterval    time.Duration
	restoreRedWaitTimeout     time.Duration
)

var updateRestoreRedCmd = &cobra.Command{
	Use:   "restore-red",
	Short: "Restore red (unassigned-primary) indices from a snapshot, in batches",
	Long: utils.Trim(`
Recover indices whose primaries are unassigned (health red) by restoring them from a snapshot.

The command intersects the indices contained in the snapshot with the red indices matching
--pattern, so only recoverable, currently-broken indices are touched. Matches are restored in
batches: each batch is closed first (a closed index can be restored over), then the restore is
submitted asynchronously and this command polls until the batch's indices are open and no longer
red before moving on. Requests stay short, so the restore is not affected by client or proxy
timeouts on long-running calls.

Indices that are closed and present in the snapshot are reported but skipped by default; pass
--include-closed to restore them too (an interrupted earlier run can leave indices closed).

This is the snapshot-based counterpart to the shard-recovery flow (get shard-stores / update
reroute): use it when no in-cluster copy survives but the data exists in a snapshot.

Instead of selecting red indices with --pattern, --alias-pattern selects every index currently
backing a matching alias (e.g. "logz-*-write-alias"), whatever its health — the disaster-recovery
flow where active indices are re-seeded from the last snapshot. In this mode closed indices are
restored too, and --exclude-today-tomorrow skips indices still being written (names carrying
today's or tomorrow's yymmdd date). --restore-replicas, --box-type and --ignore-index-setting
override index settings on the restored indices.`),
	Example: utils.TrimAndIndent(`
# Preview what would be restored
esctl update restore-red --repository my-repo --snapshot snap-1 --pattern "logz-*" --dry-run

# Restore red indices, renaming write-aliases so they don't collide with live ones
esctl update restore-red --repository my-repo --snapshot snap-1 --pattern "logz-*" \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias"

# Also pick up indices left closed by an interrupted earlier run
esctl update restore-red --repository my-repo --snapshot snap-1 --pattern "logz-*" --include-closed

# DR flow: re-seed all write-alias indices from the snapshot, move their aliases out
# of the way, drop replicas, and re-pin allocation to the default/ingestion tier
esctl update restore-red --repository my-repo --snapshot snap-1 \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateSelectionFlags(restoreRedPattern, restoreRedAliasPattern); err != nil {
			return err
		}
		aliasMode := restoreRedAliasPattern != ""
		catPattern := restoreRedPattern
		if aliasMode {
			catPattern = restoreRedAliasPattern
		}

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

		catClient := cat.NewCat()
		indices, err := catClient.CatIndices(ctx, "", catPattern, "")
		if err != nil {
			return fmt.Errorf("failed to list indices for pattern %q: %w", catPattern, err)
		}

		var exclusions []string
		if restoreRedExcludeToday {
			exclusions = dateExclusions(time.Now())
		}
		sel := selectIndices(indices, inSnapshot, aliasMode, restoreRedIncludeClosed, exclusions)

		if len(sel.dateExcluded) > 0 {
			fmt.Printf("note: %d index(es) skipped: their names carry today's or tomorrow's date, so they are still being written.\n",
				len(sel.dateExcluded))
		}
		if len(sel.notInSnapshot) > 0 {
			fmt.Printf("note: %d index(es) are not in the snapshot and cannot be restored from it: %s\n",
				len(sel.notInSnapshot), strings.Join(sel.notInSnapshot, ", "))
		}
		if len(sel.closedSkipped) > 0 {
			fmt.Printf("note: %d matching index(es) in the snapshot are closed and were skipped.\n"+
				"      An interrupted earlier run can leave indices closed; re-run with --include-closed to restore them.\n",
				len(sel.closedSkipped))
		}
		toRestore := sel.toRestore
		if len(toRestore) == 0 {
			fmt.Println("Nothing to restore: no matching indices are present in the snapshot.")
			return nil
		}
		sort.Sort(sort.Reverse(sort.StringSlice(toRestore)))

		if restoreRedIncludeAliases && restoreRedRenameAliasPat == "" {
			fmt.Println("warning: aliases are restored as-is (--include-aliases is on and no --rename-alias-pattern was given);\n" +
				"         a restored write-alias can collide with the alias of a live index. Pass --rename-alias-pattern/" +
				"--rename-alias-replacement,\n         or --include-aliases=false, to avoid this.")
		}

		batches := batchStrings(toRestore, restoreRedBatchSize)
		fmt.Printf("Found %d index(es) to restore from snapshot %q in %d batch(es) of up to %d.\n",
			len(toRestore), restoreRedSnapshot, len(batches), restoreRedBatchSize)

		idxClient := index.NewIndex()
		for n, b := range batches {
			label := fmt.Sprintf("batch %d/%d", n+1, len(batches))
			joined := strings.Join(b, ",")

			if restoreRedDryRun {
				fmt.Printf("[dry-run] %s (%d indices): would %srestore %s\n",
					label, len(b), closeVerb(restoreRedClose), joined)
				continue
			}

			fmt.Printf("%s: %d index(es)\n", label, len(b))
			if restoreRedClose {
				resp, err := idxClient.Close(ctx, b)
				if err != nil {
					return fmt.Errorf("%s: failed to close indices: %w", label, err)
				}
				if !resp.Acknowledged {
					return fmt.Errorf("%s: close was not acknowledged by the cluster; aborting before restore", label)
				}
			}

			req := buildRestoreRequest(joined, restoreRedIncludeAliases,
				restoreRedRenameAliasPat, restoreRedRenameAliasRepl,
				restoreRedReplicas, restoreRedBoxType, restoreRedIgnoreSettings)
			// Submitted asynchronously: the request returns as soon as the cluster
			// accepts the restore, keeping the HTTP call short. Progress is then
			// observed by polling, which survives client and proxy timeouts.
			if err := snapshots.RestoreSnapshotWithTimeout(ctx, restoreRedRepo, restoreRedSnapshot, req, false, restoreRedCMTimeout); err != nil {
				return fmt.Errorf("%s: failed to restore: %w%s", label, err, alreadyExistsHint(err))
			}

			if !restoreRedWait {
				continue
			}
			if err := waitForBatch(ctx, catClient, pollPatternFor(aliasMode, restoreRedPattern), b, label); err != nil {
				return fmt.Errorf("%s: %w", label, err)
			}
		}

		if restoreRedDryRun {
			fmt.Printf("Dry-run complete. Would restore %d index(es).\n", len(toRestore))
		} else {
			fmt.Printf("Done. Restored %d index(es).\n", len(toRestore))
		}
		return nil
	},
}

// waitForBatch polls until every index in the batch is open and no longer red,
// which is the point at which its primaries have been restored and assigned.
func waitForBatch(ctx context.Context, catClient cat.Cat, pattern string, batch []string, label string) error {
	want := make(map[string]bool, len(batch))
	for _, idx := range batch {
		want[idx] = true
	}

	deadline := time.Now().Add(restoreRedWaitTimeout)
	lastPending := -1
	for {
		// Query by pattern rather than by name: a pattern tolerates indices that
		// briefly do not exist, where an exact-name lookup would 404.
		indices, err := catClient.CatIndices(ctx, "", pattern, "")
		if err != nil {
			return fmt.Errorf("failed to poll index status: %w", err)
		}

		pending := 0
		for _, i := range indices {
			if want[i.Index] && (isClosedIndex(i) || strings.EqualFold(i.Health, "red")) {
				pending++
			}
		}
		if pending == 0 {
			fmt.Printf("  %s: restore complete\n", label)
			return nil
		}
		if pending != lastPending {
			fmt.Printf("  %s: %d/%d index(es) still restoring...\n", label, pending, len(batch))
			lastPending = pending
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s with %d index(es) still not recovered; "+
				"inspect with 'esctl get recovery' or 'esctl get shard-stores --status red'",
				restoreRedWaitTimeout, pending)
		}

		select {
		case <-time.After(restoreRedPollInterval):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// isClosedIndex reports whether _cat/indices shows the index as closed.
func isClosedIndex(i cat.CatIndiceResponse) bool {
	return strings.EqualFold(i.Status, "close")
}

// alreadyExistsHint explains the most common restore failure: the target index is
// still open, so it must be closed (or deleted) before it can be restored over.
func alreadyExistsHint(err error) string {
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		return ""
	}
	return "\nhint: the index is still open. Restore requires it closed — re-run with --close (the default), " +
		"or use --rename-alias-pattern/--rename-alias-replacement to restore under different alias names."
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

// dateExclusions returns the yymmdd stamps for now and the next calendar day.
// Indices carrying these dates are still being written and must not be
// restored over.
func dateExclusions(now time.Time) []string {
	return []string{now.Format("060102"), now.AddDate(0, 0, 1).Format("060102")}
}

// containsAny reports whether s contains any of the substrings.
func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// selection classifies the candidate indices for a restore run.
type selection struct {
	toRestore     []string
	closedSkipped []string
	notInSnapshot []string
	dateExcluded  []string
}

// selectIndices classifies candidates for restore. Alias mode takes every
// index regardless of health or open/closed state; pattern mode keeps the
// original semantics (red indices, closed ones only with includeClosed).
// Indices whose names carry an excluded date stamp are set aside first, and
// candidates absent from the snapshot are reported rather than restored.
func selectIndices(indices []cat.CatIndiceResponse, inSnapshot map[string]bool, aliasMode, includeClosed bool, exclusions []string) selection {
	var sel selection
	for _, i := range indices {
		if containsAny(i.Index, exclusions) {
			sel.dateExcluded = append(sel.dateExcluded, i.Index)
			continue
		}

		candidate := false
		switch {
		case aliasMode:
			candidate = true
		case isClosedIndex(i):
			// A closed index is not serving traffic. It is often the residue
			// of an interrupted restore, but it can also be closed
			// deliberately, so it is only restored when explicitly requested.
			if includeClosed {
				candidate = true
			} else if inSnapshot[i.Index] {
				sel.closedSkipped = append(sel.closedSkipped, i.Index)
			}
		case strings.EqualFold(i.Health, "red"):
			candidate = true
		}
		if !candidate {
			continue
		}

		if !inSnapshot[i.Index] {
			sel.notInSnapshot = append(sel.notInSnapshot, i.Index)
			continue
		}
		sel.toRestore = append(sel.toRestore, i.Index)
	}
	return sel
}

// buildRestoreRequest assembles the restore body for one batch. replicas < 0
// and an empty boxType leave the corresponding snapshot settings untouched.
func buildRestoreRequest(joined string, includeAliases bool, renamePat, renameRepl string, replicas int, boxType string, ignoreSettings []string) snapshots.RestoreSnapshotRequest {
	req := snapshots.RestoreSnapshotRequest{
		Indices:                joined,
		IncludeAliases:         includeAliases,
		RenameAliasPattern:     renamePat,
		RenameAliasReplacement: renameRepl,
		IgnoreIndexSettings:    ignoreSettings,
	}
	settings := map[string]any{}
	if replicas >= 0 {
		settings["index.number_of_replicas"] = replicas
	}
	if boxType != "" {
		settings["index.routing.allocation.include.box_type"] = boxType
	}
	if len(settings) > 0 {
		req.IndexSettings = settings
	}
	return req
}

// pollPatternFor picks the _cat/indices pattern used while waiting for a
// batch. Alias mode cannot poll by the alias pattern: once an index is
// restored with its alias renamed it no longer matches and would look
// complete while its shards are still restoring, so poll everything and let
// waitForBatch filter by name.
func pollPatternFor(aliasMode bool, pattern string) string {
	if aliasMode {
		return "*"
	}
	return pattern
}

// validateSelectionFlags enforces that exactly one selection mode is chosen.
func validateSelectionFlags(pattern, aliasPattern string) error {
	if (pattern == "") == (aliasPattern == "") {
		return fmt.Errorf("exactly one of --pattern or --alias-pattern must be provided")
	}
	return nil
}

func init() {
	updateRestoreRedCmd.Flags().StringVar(&restoreRedRepo, "repository", "", "Snapshot repository (required)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedSnapshot, "snapshot", "", "Snapshot name (required)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedPattern, "pattern", "", "Index pattern to match red indices (exactly one of --pattern/--alias-pattern is required)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedAliasPattern, "alias-pattern", "", "Select every index backing an alias that matches this pattern, regardless of health (exactly one of --pattern/--alias-pattern is required)")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedExcludeToday, "exclude-today-tomorrow", false, "Skip indices whose names contain today's or tomorrow's date (yymmdd, local time)")
	updateRestoreRedCmd.Flags().IntVar(&restoreRedReplicas, "restore-replicas", -1, "Override index.number_of_replicas on restored indices (-1 keeps the snapshot value)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedBoxType, "box-type", "", "Set index.routing.allocation.include.box_type on restored indices (e.g. \"default,ingestion\")")
	updateRestoreRedCmd.Flags().StringArrayVar(&restoreRedIgnoreSettings, "ignore-index-setting", nil, "Index setting to drop from the snapshot on restore (repeatable), e.g. index.routing.allocation.total_shards_per_node")
	updateRestoreRedCmd.Flags().IntVar(&restoreRedBatchSize, "batch-size", 50, "Number of indices to restore per batch")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedClose, "close", true, "Close each index before restoring it")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedIncludeClosed, "include-closed", false, "Also restore matching indices that are currently closed")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedIncludeAliases, "include-aliases", true, "Restore the aliases stored in the snapshot")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedRenameAliasPat, "rename-alias-pattern", "", "Regex to match alias names to rename on restore")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedRenameAliasRepl, "rename-alias-replacement", "", "Replacement for renamed aliases (may reference $1, $2, ...)")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedWait, "wait", true, "Wait for each batch to finish restoring before starting the next")
	updateRestoreRedCmd.Flags().DurationVar(&restoreRedPollInterval, "poll-interval", 10*time.Second, "How often to poll index status while waiting")
	updateRestoreRedCmd.Flags().DurationVar(&restoreRedWaitTimeout, "wait-timeout", 30*time.Minute, "Maximum time to wait for a single batch to recover")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedCMTimeout, "cluster-manager-timeout", "5m", "Cluster-manager (master) timeout per restore request")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedDryRun, "dry-run", false, "Show what would be restored without changing anything")
	_ = updateRestoreRedCmd.MarkFlagRequired("repository")
	_ = updateRestoreRedCmd.MarkFlagRequired("snapshot")

	updateCmd.AddCommand(updateRestoreRedCmd)
}
