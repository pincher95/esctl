# esctl Command Structure

This document reflects the current `esctl` command structure.

## Core verbs

`esctl` uses verb-first commands for resource operations:

- `get` — read or list resources
- `set` — create or update resources
- `delete` — remove resources
- `update` — operational updates (reroute, cache clear, alias move, snapshot restore, pipeline simulate, etc.)

## Top-level commands

These remain as dedicated command groups:

- `describe` — detailed JSON/YAML views (cluster, index, node)
- `count` — document counting with filters and grouping
- `query` — query execution
- `bulk` — bulk operations
- `analyze` — text analysis
- `explain` — query match explanation
- `profile` — query profiling
- `config` — configuration management
- `version` — version info

> ILM, index templates, index maintenance, and task cancellation are no longer
> dedicated noun trees — they were migrated to verb-first commands (see below).
> The old `esctl ilm`, `esctl template`, `esctl index`, and `esctl tasks` commands
> have been removed.

## Resource commands

### Aliases
```bash
esctl get aliases                                          # List all aliases
esctl get aliases --name <alias>                           # Get specific alias
esctl set alias --name <alias> --indices <idx1,idx2>       # Add alias
esctl delete alias --name <alias> --indices <idx1,idx2>    # Remove alias
esctl update alias --name <alias> --from <idx> --to <idx>  # Move alias
```

### Pipelines
```bash
esctl get pipelines                                        # List all pipelines
esctl get pipelines --id <pipeline-id>                     # Get specific pipeline
esctl set pipeline --id <pipeline-id> --file pipeline.json # Create/update
esctl delete pipeline --id <pipeline-id>                   # Delete
esctl update pipeline --file request.json                  # Simulate
```

### Snapshots
```bash
esctl get snapshot --repository <repo>                     # List snapshots in repo
esctl get snapshot --repository <repo> --name <snapshot>   # Get specific snapshot
esctl get snapshot-status                                  # Snapshot status
esctl get snapshot-repos                                   # List repositories
esctl set snapshot --repository <repo> --name <snapshot>   # Create snapshot
esctl update snapshot --repository <repo> --name <snapshot> # Restore
esctl delete snapshot --repository <repo> --name <snapshot> # Delete
esctl set snapshot-repo --repository <repo> --type fs --settings "location:/backup"
esctl delete snapshot-repo --repository <repo>
```

### Security (users/roles)
```bash
esctl get users                                            # List users
esctl get users --name <username>                          # Get user
esctl set user --name <username> --file user.json          # Create/update
esctl delete user --name <username>                        # Delete
esctl get roles                                            # List roles
esctl get roles --name <role>                              # Get role
esctl set role --name <role> --file role.json              # Create/update
esctl delete role --name <role>                            # Delete
```

### Reindex
```bash
esctl set reindex --source <index> --dest <index>          # Start reindex
esctl get reindex --task-id <task-id>                      # Get status
esctl delete reindex --task-id <task-id>                   # Cancel
```

### Scripts
```bash
esctl get scripts                                          # List all scripts
esctl get scripts --id <script-id>                         # Get specific script
esctl set script --id <id> --lang painless --source "..."  # Create/update
esctl delete script --id <script-id>                       # Delete
```

### Search Templates
```bash
esctl get search-templates                                 # List all
esctl get search-templates --id <template-id>              # Get specific
esctl set search-template --id <id> --file template.json   # Create/update
esctl delete search-template --id <template-id>            # Delete
esctl update search-template render --id <id> --params '{}' # Render
```

### Data Streams
```bash
esctl get data-streams                                     # List all
esctl get data-streams --name <name>                       # Get specific
esctl delete data-stream --name <name>                     # Delete
esctl update data-stream rollover --name <name>            # Rollover
```

### ILM (Index Lifecycle Management)
Verb-first commands. The old `esctl ilm <verb>` tree has been removed.
```bash
esctl get ilm-policies                                     # List all policies
esctl get ilm-policies --name <policy>                     # Get a specific policy
esctl set ilm-policy --name <policy> --file policy.json    # Create/update
esctl delete ilm-policy --name <policy> [--force]          # Delete
esctl get ilm-explain --index <index>                      # Explain ILM status for indices
esctl update ilm-retry --index <index>                     # Retry a failed ILM step
```

Migration:
- `esctl ilm list`            → `esctl get ilm-policies`
- `esctl ilm get --name`      → `esctl get ilm-policies --name`
- `esctl ilm put`             → `esctl set ilm-policy`
- `esctl ilm delete`          → `esctl delete ilm-policy`
- `esctl ilm explain`         → `esctl get ilm-explain`
- `esctl ilm retry`           → `esctl update ilm-retry`

### SLM (Snapshot Lifecycle Management)
```bash
esctl get slm-policies                                     # List all SLM policies
esctl get slm-policies --name <policy>                     # Get a specific policy
esctl set slm-policy --name <policy> --file policy.json    # Create/update
esctl delete slm-policy --name <policy> [--force]          # Delete
esctl update slm-execute --name <policy>                   # Take a snapshot now (off-schedule)
```

### Index templates
Verb-first commands. The old `esctl template ...` and `esctl template component ...` trees
have been removed.
```bash
esctl get templates                                        # List index templates (composable + legacy)
esctl get templates --name <template>                      # Get a specific template
esctl set template --name <t> --file template.json         # Create/update from file
esctl set template --name <t> --patterns "logs-*" --priority 100
esctl delete template --name <template> [--force]          # Delete

esctl get component-templates                              # List component templates
esctl get component-templates --name <name>                # Get a specific component template
esctl set component-template --name <name> --file c.json   # Create/update
esctl delete component-template --name <name> [--force]    # Delete
```

Migration:
- `esctl template list/get`               → `esctl get templates [--name]`
- `esctl template put/delete`             → `esctl set/delete template`
- `esctl template component list/get`     → `esctl get component-templates [--name]`
- `esctl template component put/delete`   → `esctl set/delete component-template`

### Tasks
```bash
esctl get tasks                                            # List running tasks
esctl delete task --task-id <id>                           # Cancel a task (replaces 'esctl tasks cancel')
esctl delete task --actions "*reindex"                     # Cancel tasks by action
```

### Index maintenance (moved under `update index`)
The standalone `esctl index refresh/flush/forcemerge` tree has been removed; all index actions now
live under `update index`.
```bash
esctl update index refresh --index <idx>                   # was: esctl index refresh
esctl update index flush --index <idx>                     # was: esctl index flush
esctl update index forcemerge --index <idx> --max-num-segments 1
```

### Cluster and Diagnostics
```bash
esctl get health                                           # Cluster health
esctl get nodes                                            # List nodes
esctl get indices                                          # List indices
esctl get shards                                           # List shards
esctl get allocation                                       # Allocation info
esctl get shard-stores                                     # Shard copies on disk (diagnose unassigned shards)
esctl get shard-stores --status red                        # Only shards with an unassigned primary
esctl get fielddata                                        # Fielddata cache
esctl get plugins                                          # Installed plugins
esctl get tasks                                            # Running tasks
esctl get thread-pools                                     # Thread pool stats
esctl get hot-threads                                      # Hot threads
esctl get node-stats                                       # Node health: heap, GC, thread-pool rejections, breakers, disk
esctl get pending-tasks                                    # Queued cluster-state changes (diagnose a stuck master)
esctl get index-stats --indices <index>                    # Index statistics
esctl get recovery --indices <index>                       # Recovery status
esctl get segments --indices <index>                       # Segment info
```

### Shard allocation recovery
Diagnose, then act on unassigned shards:
```bash
esctl get shards --unassigned                              # What is unassigned
esctl get explain                                          # Why (allocation deciders)
esctl get shard-stores --status red                        # Which nodes still hold a copy on disk
esctl update reroute retry-failed                          # Retry allocation blocked by prior failures

# Manual allocation commands (destructive ones require --accept-data-loss):
esctl update reroute move --index <idx> --shard 0 --from-node <a> --to-node <b>
esctl update reroute cancel --index <idx> --shard 0 --node <n> [--allow-primary]
esctl update reroute allocate-stale-primary --index <idx> --shard 0 --node <n> --accept-data-loss
esctl update reroute allocate-empty-primary --index <idx> --shard 0 --node <n> --accept-data-loss

# When no in-cluster copy survives, restore the red indices from a snapshot. Each batch is
# closed, the restore is submitted asynchronously, then the batch is polled until its indices
# are open and no longer red (so client/proxy timeouts cannot break a long restore).
esctl update restore-red --repository <repo> --snapshot <snap> --pattern "logz-*" [--dry-run]

# Rename restored aliases so they cannot collide with live write-aliases:
esctl update restore-red --repository <repo> --snapshot <snap> --pattern "logz-*" \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-\$1-alias"

# Pick up indices left closed by an interrupted earlier run:
esctl update restore-red --repository <repo> --snapshot <snap> --pattern "logz-*" --include-closed

# DR flow: re-seed all write-alias indices from the snapshot with settings overrides
esctl update restore-red --repository <repo> --snapshot <snap> \
  --alias-pattern "logz-*-write-alias" --exclude-today-tomorrow \
  --rename-alias-pattern "logz-(.+)-write-alias" --rename-alias-replacement "old-$1-alias" \
  --restore-replicas 0 --box-type "default,ingestion" \
  --ignore-index-setting index.routing.allocation.total_shards_per_node \
  --ignore-index-setting index.routing.allocation.require._ip
```

## Watch mode

Every `get` command honors the shared `--watch/-w` and `--interval` flags (default 5s) and
redraws in place. Press `q` (or Esc) to quit immediately, or Ctrl-C. The terminal is always
restored on exit. (`get hot-threads` is the one exception — it has its own sampling `--interval`.)

## Global flags

```
--host string        Elasticsearch host
--port int           Elasticsearch port (default 9200)
--protocol string    Elasticsearch protocol (default "http")
--username string    Elasticsearch username
--password string    Elasticsearch password
--api-key string     API key (Authorization: ApiKey ...); takes precedence over username/password ($ESCTL_API_KEY)
--ca-cert string     Path to a PEM CA bundle to verify the server's TLS certificate ($ESCTL_CA_CERT)
--insecure           Skip TLS verification (INSECURE; testing only — prefer --ca-cert)
--context string     Override context from config file
--timeout duration   Global timeout (e.g. 30s, 2m)
--debug              Enable debug mode
-o, --output string  Output format: table|json|yaml (default "table")
```

## Shell completion

`esctl` ships Cobra-generated completion. Enable it once per shell, e.g.:

```bash
# zsh
esctl completion zsh > "${fpath[1]}/_esctl"    # then restart your shell
# bash
esctl completion bash | sudo tee /etc/bash_completion.d/esctl >/dev/null
# fish
esctl completion fish > ~/.config/fish/completions/esctl.fish
```
