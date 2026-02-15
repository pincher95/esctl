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

- `template` — index template management (list, get, put, delete, exists; includes component subcommand)
- `ilm` — ILM policy management (list, get, put, delete, exists, explain, retry)
- `index` — index maintenance (refresh, flush, forcemerge)
- `describe` — detailed JSON/YAML views (cluster, index)
- `count` — document counting with filters and grouping
- `query` — query execution
- `bulk` — bulk operations
- `analyze` — text analysis
- `explain` — query match explanation
- `profile` — query profiling
- `tasks` — task management (cancel)
- `config` — configuration management
- `version` — version info

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

### Cluster and Diagnostics
```bash
esctl get health                                           # Cluster health
esctl get nodes                                            # List nodes
esctl get indices                                          # List indices
esctl get shards                                           # List shards
esctl get allocation                                       # Allocation info
esctl get fielddata                                        # Fielddata cache
esctl get plugins                                          # Installed plugins
esctl get tasks                                            # Running tasks
esctl get thread-pools                                     # Thread pool stats
esctl get hot-threads                                      # Hot threads
esctl get index-stats --indices <index>                    # Index statistics
esctl get recovery --indices <index>                       # Recovery status
esctl get segments --indices <index>                       # Segment info
```

## Global flags

```
--host string        Elasticsearch host
--port int           Elasticsearch port (default 9200)
--protocol string    Elasticsearch protocol (default "http")
--username string    Elasticsearch username
--password string    Elasticsearch password
--context string     Override context from config file
--timeout duration   Global timeout (e.g. 30s, 2m)
--debug              Enable debug mode
-o, --output string  Output format: table|json|yaml (default "table")
```
