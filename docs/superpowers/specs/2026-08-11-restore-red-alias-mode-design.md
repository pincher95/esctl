# Design: `update restore-red` — alias-based selection, date exclusion, and index-settings overrides

**Date:** 2026-08-11
**Status:** Approved

## Background

A dev-provided DR bash script restores the *active* indices of a Logz.io production
cluster from a snapshot. Its behavior differs from the existing
`esctl update restore-red` command in three ways:

1. It selects indices by resolving the write-alias pattern
   `logz-*-write-alias` via `_cat/indices` — taking every backing index
   regardless of health — instead of selecting only red indices.
2. It excludes indices whose names contain today's or tomorrow's date
   (`yymmdd`), because those are still being written.
3. Its restore body overrides index settings: `number_of_replicas: 0`,
   clears `routing.allocation.total_shards_per_node` and
   `routing.allocation.require._ip`, and sets
   `routing.allocation.include.box_type: "default,ingestion"`.

The script itself is buggy (its `ACTIVE_INDICES`/`ALL_INDICES` intersection is
a no-op, and it POSTs to a nonexistent `/{indices}/{repo}/_settings` endpoint
with a restore body). This design implements the script's *intent* — the
snapshot `_restore` API — on top of the existing command.

**Decision:** extend `update restore-red` rather than add a sibling command.
~90% of the machinery (snapshot intersection, batching, close-before-restore,
async submit + poll, alias rename, dry-run) is shared.

## CLI surface

New flags on `esctl update restore-red`:

| Flag | Type / default | Effect |
|---|---|---|
| `--alias-pattern` | string, `""` | Select all indices currently backing aliases matching this pattern, any health. Exactly one of `--pattern` / `--alias-pattern` is required. |
| `--exclude-today-tomorrow` | bool, `false` | Skip candidates whose name contains today's or tomorrow's `yymmdd` (local time). Works in both selection modes. |
| `--restore-replicas` | int, `-1` | When ≥ 0, sets `index_settings["index.number_of_replicas"]` on the restore request. |
| `--box-type` | string, `""` | When non-empty, sets `index_settings["index.routing.allocation.include.box_type"]`. |
| `--ignore-index-setting` | string array, empty | Each value appended to the restore request's `ignore_index_settings`. |

Existing flags are unchanged. `--pattern` loses its `MarkFlagRequired` in
favor of a `PreRunE`/`RunE` check that exactly one of `--pattern` or
`--alias-pattern` was provided.

### Canonical DR invocation (added to command examples)

```
esctl update restore-red --repository my-repo --snapshot snap-1 \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip
```

`ignore_index_settings` is the API-correct translation of the script's
`"total_shards_per_node": null` / `"require": {"_ip": null}`: restore's
`index_settings` does not reliably accept nulls; `ignore_index_settings`
drops the snapshot-stored setting.

The box_type settings path is `index.routing.allocation.include.box_type`,
matching the script's JSON (`"routing.allocation": {"include": {"box_type": ...}}`).

## Selection logic

Both modes share: intersect candidates with the snapshot's index list
(candidates not in the snapshot are reported and skipped), apply the date
exclusion if requested, reverse-sort (newest first), batch (default 50).

- **Pattern mode (existing, unchanged):** `_cat/indices/<pattern>`, keep red
  indices; closed ones are skipped unless `--include-closed`.
- **Alias mode (new):** `_cat/indices/<alias-pattern>` resolves aliases to
  their backing indices; keep **all** of them regardless of health.
  Already-closed indices are included automatically (restore works over a
  closed index; `--include-closed` is not consulted). Open indices are closed
  first by the existing `--close` behavior (default on).

Date exclusion: compute `time.Now()` and its calendar-day successor
(`AddDate(0, 0, 1)`, matching the script's `date -v+1d`) formatted as `060102`
(`yymmdd`), skip any candidate whose index name contains either substring.
The "now" value is injected into the filtering function so tests don't depend
on the wall clock.

## Restore request

`es/snapshots.RestoreSnapshotRequest` already has `IndexSettings
map[string]any` and `IgnoreIndexSettings []string`; no ES-layer changes.
The command builds:

- `IndexSettings["index.number_of_replicas"] = N` when `--restore-replicas ≥ 0`
- `IndexSettings["index.routing.allocation.include.box_type"] = v` when `--box-type` non-empty
- `IgnoreIndexSettings = values of --ignore-index-setting`

All other request fields (indices, include_aliases, alias rename) as today.

## Poll-after-restore fix

`waitForBatch` currently polls `_cat/indices/<pattern>` and counts pending
among the batch's names. In alias mode this is wrong: once an index is
restored with its write-alias renamed (e.g. to `old-*-alias`), it no longer
matches `logz-*-write-alias`, so an alias-pattern poll would drop it from the
result and report the batch complete while shards are still restoring.

Fix: in alias mode, `waitForBatch` polls `_cat/indices/*` and filters by the
batch's index names (the existing `want` map). Pattern mode keeps polling by
`--pattern` as today.

## Error handling

- Exactly-one-of `--pattern`/`--alias-pattern` violation → usage error before
  any cluster call.
- Everything else reuses the existing paths: unacknowledged close aborts the
  batch, restore errors get the `already exists` hint, wait timeout points at
  `get recovery` / `get shard-stores`.
- The existing warning when `--include-aliases` is on without an alias-rename
  pattern applies to both modes.

## Testing

- Table-driven unit tests for the selection/exclusion helpers: alias-mode
  vs pattern-mode filtering, closed-index handling per mode, date exclusion
  with an injected "now".
- Marshal/mock-server test asserting the restore body carries
  `index_settings` (replicas + box_type) and `ignore_index_settings`.
- Poll-target test: alias mode polls `*`, pattern mode polls the pattern.
- Manual verification: `--dry-run` against the staging cluster (creds
  provided), then a real restore on staging if feasible.

## Docs

Update the command's `Long`/`Example` text, `README.md`, and
`COMMAND_STRUCTURE.md` to describe the new flags and the DR flow.

## Out of scope

- Changes to `esctl snapshot restore` (the generic restore command).
- Replicating the script's `ACTIVE`/`ALL` grep intersection (a no-op).
- Index *name* renaming on restore (the script renames aliases only).
