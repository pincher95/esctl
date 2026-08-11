# restore-red Alias Mode Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend `esctl update restore-red` with alias-based index selection, today/tomorrow date exclusion, and index-settings overrides on the restore request, per `docs/superpowers/specs/2026-08-11-restore-red-alias-mode-design.md`.

**Architecture:** All changes live in `cmd/update/restore_red.go` (plus its test file and docs). The inline selection loop is extracted into pure, table-testable helpers; new flags feed the existing `es/snapshots.RestoreSnapshotRequest` fields (`IndexSettings`, `IgnoreIndexSettings`), which already exist — the ES layer is untouched.

**Tech Stack:** Go 1.26, cobra, stdlib `testing` (no assertion libraries — this repo uses plain `t.Errorf`/`t.Fatalf`).

## Global Constraints

- Go 1.26; run `gofmt -w` on every file you touch before committing.
- Do NOT modify anything under `es/` — `RestoreSnapshotRequest` already has `IndexSettings map[string]any` and `IgnoreIndexSettings []string` (es/snapshots/snapshots.go:93-105).
- Follow the file's existing style: package-level flag vars named `restoreRed*`, helpers at the bottom of `restore_red.go`, plain-stdlib tests in `restore_red_test.go` (package `update`, no test frameworks).
- The box_type settings key is exactly `index.routing.allocation.include.box_type` (with `.include`).
- Date stamps use Go layout `060102` (yymmdd) in local time; "tomorrow" is `AddDate(0, 0, 1)` (calendar day), never `+24h`.
- Existing pattern-mode behavior must not change except where the spec says (red/closed classification, `--include-closed`, notes, batching, polling by `--pattern`).
- All commands below run from the repo root `/Users/yuri.tsuprun/Documents/GitHub/esctl`.

---

### Task 1: Date-exclusion helper

**Files:**
- Modify: `cmd/update/restore_red.go` (add helpers at bottom, near `batchStrings`)
- Test: `cmd/update/restore_red_test.go`

**Interfaces:**
- Consumes: nothing new.
- Produces: `dateExclusions(now time.Time) []string` — returns exactly two `yymmdd` stamps: `now` and the next calendar day. Also `containsAny(s string, subs []string) bool`. Task 2's `selectIndices` uses both; Task 4's `RunE` calls `dateExclusions(time.Now())`.

- [ ] **Step 1: Write the failing tests**

Append to `cmd/update/restore_red_test.go` (add `"reflect"` and `"time"` to the imports):

```go
func TestDateExclusions(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 0, 0, 0, time.UTC)
	if got, want := dateExclusions(now), []string{"260811", "260812"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dateExclusions(%v) = %v, want %v", now, got, want)
	}

	// Tomorrow is the next calendar day, across month and year boundaries.
	eoy := time.Date(2026, 12, 31, 23, 30, 0, 0, time.UTC)
	if got, want := dateExclusions(eoy), []string{"261231", "270101"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dateExclusions(%v) = %v, want %v", eoy, got, want)
	}
}

func TestContainsAny(t *testing.T) {
	if !containsAny("logz-abc-260811-000001", []string{"260811", "260812"}) {
		t.Error("expected match on 260811")
	}
	if containsAny("logz-abc-260810-000001", []string{"260811", "260812"}) {
		t.Error("did not expect a match for 260810")
	}
	if containsAny("logz-abc-260811", nil) {
		t.Error("nil exclusions must never match")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/update/ -run 'TestDateExclusions|TestContainsAny' -v`
Expected: FAIL to compile with `undefined: dateExclusions` / `undefined: containsAny`

- [ ] **Step 3: Write the implementation**

Append to `cmd/update/restore_red.go` (after `batchStrings`):

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/update/ -run 'TestDateExclusions|TestContainsAny' -v`
Expected: PASS (2 tests)

- [ ] **Step 5: Commit**

```bash
gofmt -w cmd/update/restore_red.go cmd/update/restore_red_test.go
git add cmd/update/restore_red.go cmd/update/restore_red_test.go
git commit -m "feat(restore-red): add date-exclusion helpers for today/tomorrow indices"
```

---

### Task 2: Selection classifier (`selectIndices`)

**Files:**
- Modify: `cmd/update/restore_red.go`
- Test: `cmd/update/restore_red_test.go`

**Interfaces:**
- Consumes: `containsAny` (Task 1), existing `isClosedIndex(i cat.CatIndiceResponse) bool` (restore_red.go:225), `cat.CatIndiceResponse` (fields used: `Index`, `Health`, `Status` — es/cat/indices.go:11).
- Produces:

```go
type selection struct {
	toRestore     []string
	closedSkipped []string
	notInSnapshot []string
	dateExcluded  []string
}

func selectIndices(indices []cat.CatIndiceResponse, inSnapshot map[string]bool, aliasMode, includeClosed bool, exclusions []string) selection
```

Task 4's `RunE` replaces its inline classification loop with this call.

Classification rules (this is the heart of the feature — implement exactly):
1. Name contains any exclusion stamp → `dateExcluded`, nothing else considered.
2. Alias mode: every remaining index is a candidate regardless of health or open/closed state (`--include-closed` is NOT consulted).
3. Pattern mode: closed indices are candidates only with `includeClosed` (otherwise, if in the snapshot, they're recorded in `closedSkipped` as today); open indices are candidates only when health is red; anything else is silently ignored.
4. A candidate absent from the snapshot → `notInSnapshot`; present → `toRestore`.

- [ ] **Step 1: Write the failing tests**

Append to `cmd/update/restore_red_test.go` (add `"github.com/pincher95/esctl/es/cat"` to imports):

```go
func catIdx(name, health, status string) cat.CatIndiceResponse {
	return cat.CatIndiceResponse{Index: name, Health: health, Status: status}
}

func TestSelectIndicesPatternMode(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-red", "red", "open"),
		catIdx("logz-green", "green", "open"),
		catIdx("logz-closed", "red", "close"),
		catIdx("logz-red-nosnap", "red", "open"),
	}
	inSnap := map[string]bool{"logz-red": true, "logz-closed": true}

	sel := selectIndices(indices, inSnap, false, false, nil)
	if want := []string{"logz-red"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if want := []string{"logz-closed"}; !reflect.DeepEqual(sel.closedSkipped, want) {
		t.Errorf("closedSkipped = %v, want %v", sel.closedSkipped, want)
	}
	if want := []string{"logz-red-nosnap"}; !reflect.DeepEqual(sel.notInSnapshot, want) {
		t.Errorf("notInSnapshot = %v, want %v", sel.notInSnapshot, want)
	}

	// includeClosed promotes the closed index to a restore target.
	sel = selectIndices(indices, inSnap, false, true, nil)
	if want := []string{"logz-red", "logz-closed"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore with includeClosed = %v, want %v", sel.toRestore, want)
	}
	if len(sel.closedSkipped) != 0 {
		t.Errorf("closedSkipped with includeClosed = %v, want empty", sel.closedSkipped)
	}
}

func TestSelectIndicesAliasMode(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-green", "green", "open"),
		catIdx("logz-red", "red", "open"),
		catIdx("logz-closed", "green", "close"),
		catIdx("logz-nosnap", "green", "open"),
	}
	inSnap := map[string]bool{"logz-green": true, "logz-red": true, "logz-closed": true}

	// includeClosed=false on purpose: alias mode must take closed indices anyway.
	sel := selectIndices(indices, inSnap, true, false, nil)
	if want := []string{"logz-green", "logz-red", "logz-closed"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if want := []string{"logz-nosnap"}; !reflect.DeepEqual(sel.notInSnapshot, want) {
		t.Errorf("notInSnapshot = %v, want %v", sel.notInSnapshot, want)
	}
	if len(sel.closedSkipped) != 0 || len(sel.dateExcluded) != 0 {
		t.Errorf("unexpected closedSkipped=%v dateExcluded=%v", sel.closedSkipped, sel.dateExcluded)
	}
}

func TestSelectIndicesDateExclusion(t *testing.T) {
	indices := []cat.CatIndiceResponse{
		catIdx("logz-a-260811-000001", "red", "open"),
		catIdx("logz-a-260812-000001", "red", "open"),
		catIdx("logz-a-260810-000001", "red", "open"),
	}
	inSnap := map[string]bool{
		"logz-a-260811-000001": true,
		"logz-a-260812-000001": true,
		"logz-a-260810-000001": true,
	}
	sel := selectIndices(indices, inSnap, true, false, []string{"260811", "260812"})
	if want := []string{"logz-a-260810-000001"}; !reflect.DeepEqual(sel.toRestore, want) {
		t.Errorf("toRestore = %v, want %v", sel.toRestore, want)
	}
	if want := []string{"logz-a-260811-000001", "logz-a-260812-000001"}; !reflect.DeepEqual(sel.dateExcluded, want) {
		t.Errorf("dateExcluded = %v, want %v", sel.dateExcluded, want)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/update/ -run 'TestSelectIndices' -v`
Expected: FAIL to compile with `undefined: selectIndices`

- [ ] **Step 3: Write the implementation**

Append to `cmd/update/restore_red.go`:

```go
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
```

Note: this moves the "closed index" comment from `RunE` into the helper. Task 4 deletes the original inline loop.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/update/ -run 'TestSelectIndices' -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Commit**

```bash
gofmt -w cmd/update/restore_red.go cmd/update/restore_red_test.go
git add cmd/update/restore_red.go cmd/update/restore_red_test.go
git commit -m "feat(restore-red): extract selection classifier with alias mode and date exclusion"
```

---

### Task 3: Request builder, poll pattern, flag validation

**Files:**
- Modify: `cmd/update/restore_red.go`
- Test: `cmd/update/restore_red_test.go`

**Interfaces:**
- Consumes: `snapshots.RestoreSnapshotRequest` (es/snapshots/snapshots.go:93 — fields `Indices`, `IncludeAliases`, `RenameAliasPattern`, `RenameAliasReplacement`, `IndexSettings map[string]any`, `IgnoreIndexSettings []string`).
- Produces (all used by Task 4's `RunE`):

```go
func buildRestoreRequest(joined string, includeAliases bool, renamePat, renameRepl string, replicas int, boxType string, ignoreSettings []string) snapshots.RestoreSnapshotRequest
func pollPatternFor(aliasMode bool, pattern string) string
func validateSelectionFlags(pattern, aliasPattern string) error
```

- [ ] **Step 1: Write the failing tests**

Append to `cmd/update/restore_red_test.go` (add `"encoding/json"` and `"strings"` to imports):

```go
func TestBuildRestoreRequestSettings(t *testing.T) {
	req := buildRestoreRequest("a,b", true, "logz-(.+)-write-alias", "old-$1-alias",
		0, "default,ingestion",
		[]string{"index.routing.allocation.total_shards_per_node", "index.routing.allocation.require._ip"})
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{
		`"indices":"a,b"`,
		`"include_aliases":true`,
		`"rename_alias_pattern":"logz-(.+)-write-alias"`,
		`"rename_alias_replacement":"old-$1-alias"`,
		`"index.number_of_replicas":0`,
		`"index.routing.allocation.include.box_type":"default,ingestion"`,
		`"ignore_index_settings":["index.routing.allocation.total_shards_per_node","index.routing.allocation.require._ip"]`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("marshaled request missing %s in: %s", want, s)
		}
	}
}

func TestBuildRestoreRequestNoOverrides(t *testing.T) {
	req := buildRestoreRequest("a", true, "", "", -1, "", nil)
	if req.IndexSettings != nil {
		t.Errorf("IndexSettings = %v, want nil", req.IndexSettings)
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, banned := range []string{"index_settings", "ignore_index_settings"} {
		if strings.Contains(string(b), banned) {
			t.Errorf("marshaled request must not contain %q: %s", banned, b)
		}
	}
}

func TestPollPatternFor(t *testing.T) {
	// Alias mode polls "*": a restored index whose write-alias was renamed no
	// longer matches the alias pattern and would look complete while red.
	if got := pollPatternFor(true, "logz-*-write-alias"); got != "*" {
		t.Errorf("alias mode poll pattern = %q, want *", got)
	}
	if got := pollPatternFor(false, "logz-*"); got != "logz-*" {
		t.Errorf("pattern mode poll pattern = %q, want logz-*", got)
	}
}

func TestValidateSelectionFlags(t *testing.T) {
	if err := validateSelectionFlags("", ""); err == nil {
		t.Error("expected error when neither flag is set")
	}
	if err := validateSelectionFlags("logz-*", "logz-*-write-alias"); err == nil {
		t.Error("expected error when both flags are set")
	}
	if err := validateSelectionFlags("logz-*", ""); err != nil {
		t.Errorf("pattern only: unexpected error %v", err)
	}
	if err := validateSelectionFlags("", "logz-*-write-alias"); err != nil {
		t.Errorf("alias pattern only: unexpected error %v", err)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd/update/ -run 'TestBuildRestoreRequest|TestPollPatternFor|TestValidateSelectionFlags' -v`
Expected: FAIL to compile with `undefined: buildRestoreRequest` (and the others)

- [ ] **Step 3: Write the implementation**

Append to `cmd/update/restore_red.go`:

```go
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
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd/update/ -run 'TestBuildRestoreRequest|TestPollPatternFor|TestValidateSelectionFlags' -v`
Expected: PASS (4 tests)

- [ ] **Step 5: Commit**

```bash
gofmt -w cmd/update/restore_red.go cmd/update/restore_red_test.go
git add cmd/update/restore_red.go cmd/update/restore_red_test.go
git commit -m "feat(restore-red): restore-request builder with settings overrides, poll pattern, flag validation"
```

---

### Task 4: Wire flags and RunE

**Files:**
- Modify: `cmd/update/restore_red.go` (flag vars block at top, `RunE`, `Long`, `Example`, `init()`)

**Interfaces:**
- Consumes: `dateExclusions`, `selectIndices`, `buildRestoreRequest`, `pollPatternFor`, `validateSelectionFlags` (Tasks 1–3); existing `waitForBatch`, `batchStrings`, `closeVerb`, `alreadyExistsHint`.
- Produces: the finished command surface. New flags: `--alias-pattern` (string), `--exclude-today-tomorrow` (bool), `--restore-replicas` (int, default -1), `--box-type` (string), `--ignore-index-setting` (string array). `--pattern` is no longer `MarkFlagRequired`.

- [ ] **Step 1: Add the new flag variables**

In the `var (...)` block at the top of `cmd/update/restore_red.go`, after `restoreRedPattern`, add:

```go
	restoreRedAliasPattern    string
	restoreRedExcludeToday    bool
	restoreRedReplicas        int
	restoreRedBoxType         string
	restoreRedIgnoreSettings  []string
```

- [ ] **Step 2: Replace the selection and request code in RunE**

Replace the body of `RunE` (currently restore_red.go:63-174) with:

```go
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
```

(The structure is the original RunE with four changes: flag validation + `aliasMode`/`catPattern` at the top, the inline classification loop replaced by `selectIndices` + the three notes, `req` built by `buildRestoreRequest`, and `waitForBatch` fed by `pollPatternFor`. The "Nothing to restore" message is now mode-neutral.)

- [ ] **Step 3: Register the new flags and drop the pattern requirement**

In `init()`:
- Change the `--pattern` help text to: `"Index pattern to match red indices (exactly one of --pattern/--alias-pattern is required)"`.
- Delete the line `_ = updateRestoreRedCmd.MarkFlagRequired("pattern")`.
- Add after the `--pattern` registration:

```go
	updateRestoreRedCmd.Flags().StringVar(&restoreRedAliasPattern, "alias-pattern", "", "Select every index backing an alias that matches this pattern, regardless of health (exactly one of --pattern/--alias-pattern is required)")
	updateRestoreRedCmd.Flags().BoolVar(&restoreRedExcludeToday, "exclude-today-tomorrow", false, "Skip indices whose names contain today's or tomorrow's date (yymmdd, local time)")
	updateRestoreRedCmd.Flags().IntVar(&restoreRedReplicas, "restore-replicas", -1, "Override index.number_of_replicas on restored indices (-1 keeps the snapshot value)")
	updateRestoreRedCmd.Flags().StringVar(&restoreRedBoxType, "box-type", "", "Set index.routing.allocation.include.box_type on restored indices (e.g. \"default,ingestion\")")
	updateRestoreRedCmd.Flags().StringArrayVar(&restoreRedIgnoreSettings, "ignore-index-setting", nil, "Index setting to drop from the snapshot on restore (repeatable), e.g. index.routing.allocation.total_shards_per_node")
```

- [ ] **Step 4: Update Long and Example**

Append to the `Long` text (inside the `utils.Trim` backticks, after the "shard-recovery flow" paragraph):

```
Instead of selecting red indices with --pattern, --alias-pattern selects every index currently
backing a matching alias (e.g. "logz-*-write-alias"), whatever its health — the disaster-recovery
flow where active indices are re-seeded from the last snapshot. In this mode closed indices are
restored too, and --exclude-today-tomorrow skips indices still being written (names carrying
today's or tomorrow's yymmdd date). --restore-replicas, --box-type and --ignore-index-setting
override index settings on the restored indices.
```

Append to the `Example` text:

```
# DR flow: re-seed all write-alias indices from the snapshot, move their aliases out
# of the way, drop replicas, and re-pin allocation to the default/ingestion tier
esctl update restore-red --repository my-repo --snapshot snap-1 \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip
```

- [ ] **Step 5: Build and run the full test suite**

Run: `gofmt -w cmd/update/restore_red.go && go build ./... && go vet ./... && go test ./...`
Expected: build and vet clean; ALL tests pass (the new ones and every pre-existing test — `TestBatchStrings`, the `es/` suites, etc.)

- [ ] **Step 6: Smoke-check the CLI surface (no cluster needed)**

Run: `go run . update restore-red --repository r --snapshot s 2>&1 | head -3`
Expected: error `exactly one of --pattern or --alias-pattern must be provided`

Run: `go run . update restore-red --repository r --snapshot s --pattern "a" --alias-pattern "b" 2>&1 | head -3`
Expected: the same error

Run: `go run . update restore-red --help | grep -E "alias-pattern|exclude-today|restore-replicas|box-type|ignore-index-setting"`
Expected: all five new flags listed

- [ ] **Step 7: Commit**

```bash
git add cmd/update/restore_red.go
git commit -m "feat(restore-red): alias-pattern selection, date exclusion, index-settings overrides"
```

---

### Task 5: Documentation

**Files:**
- Modify: `README.md` (restore-red examples around line 192)
- Modify: `COMMAND_STRUCTURE.md` (restore-red section around line 214)

**Interfaces:**
- Consumes: the final flag names from Task 4.
- Produces: user-facing docs; nothing downstream.

- [ ] **Step 1: Add the DR example to README.md**

After the existing `--include-closed` example (README.md:198), add:

```
# DR flow: re-seed all write-alias indices from the snapshot with settings overrides
esctl update restore-red --repository my-repo --snapshot snap-1 \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip
```

- [ ] **Step 2: Add the same example to COMMAND_STRUCTURE.md**

After the existing restore-red examples (COMMAND_STRUCTURE.md:221), add the same block as Step 1. Match the surrounding placeholder style (`<repo>`, `<snap>`) used by that file:

```
esctl update restore-red --repository <repo> --snapshot <snap> \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip
```

- [ ] **Step 3: Commit**

```bash
git add README.md COMMAND_STRUCTURE.md
git commit -m "docs: document restore-red alias mode and settings overrides"
```

---

### Task 6: Staging verification (dry-run)

**Files:** none (verification only). Requires network access to the Logz.io staging VPN/network.

**Interfaces:**
- Consumes: the built CLI; staging context `cluster-102-staging` in `~/.config/esctl.yml` (host `es-logs-102.staging.us-east-1.internal.logz.io`).
- Produces: evidence the command works against a real cluster.

- [ ] **Step 1: Confirm connectivity**

Run: `go run . --context cluster-102-staging get health`
Expected: a health row from the staging cluster. If this fails, stop and report — the remaining steps need the cluster.

- [ ] **Step 2: Find a repository and snapshot**

Run: `go run . --context cluster-102-staging get snapshot-repos`
Then list snapshots in the newest repo (check `go run . --context cluster-102-staging get snapshot --help` for the exact flag, it takes a repository):
`go run . --context cluster-102-staging get snapshot --repository <repo-from-above>`
Pick the most recent successful snapshot name for the next step.

- [ ] **Step 3: Dry-run the DR invocation**

```bash
go run . --context cluster-102-staging update restore-red \
  --repository <repo> --snapshot <snap> \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip \
  --dry-run
```

Expected:
- No mutation happens (dry-run).
- Output lists batches of alias-backed indices; indices with today's/tomorrow's `yymmdd` in their names do NOT appear in any batch.
- If staging has indices not present in the snapshot, the `note: ... not in the snapshot` line appears.
- Ends with `Dry-run complete. Would restore N index(es).`

- [ ] **Step 4: Report results**

Paste the (trimmed) dry-run output in the summary to the user. Do NOT run a real restore on staging without the user explicitly asking.
