# Command Structure Update
 
 This document replaces the previous proposal. The CLI now uses a verb-first structure
 for resource operations. Old top-level commands for `alias`, `pipeline`, `security`,
 `snapshot`, `reindex`, and `list` have been removed.
 
 ## Current structure
 
 - `esctl get <resource>` — list or read resources
 - `esctl set <resource>` — create or update resources
 - `esctl delete <resource>` — delete resources
 - `esctl update <resource>` — operational updates (reroute, cache clear, alias move, snapshot restore, pipeline simulate, index settings)
 
 See `COMMAND_STRUCTURE.md` for the full reference and examples.
