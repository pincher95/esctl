# Command Standardization Summary
 
 The CLI has been standardized around verb-first commands:
 
 - `get` — read/list resources
 - `set` — create/update resources
 - `delete` — remove resources
 - `update` — operational updates (reroute, cache clear, alias move, snapshot restore, pipeline simulate, index settings)
 
 Top-level command trees for `alias`, `pipeline`, `security`, `snapshot`, `reindex`, and `list`
 were removed to eliminate overlap.
 
 ## Examples
 
 ```bash
 # Aliases
 esctl get aliases
 esctl get alias <alias>
 esctl set alias <alias> --indices=<index1,index2,...>
 esctl delete alias <alias> --indices=<index1,index2,...>
 esctl update alias <alias> --from=<from-index> --to=<to-index>
 
 # Pipelines
 esctl get pipelines
 esctl get pipeline <pipeline-id>
 esctl set pipeline <pipeline-id> --file=pipeline.json
 esctl delete pipeline <pipeline-id>
 esctl update pipeline --file=request.json [--pipeline=<pipeline-id>]
 
 # Snapshots
 esctl get snapshot <repo> <snapshot>
 esctl get snapshot <repo>
 esctl get snapshot-status [repo] [snapshot]
 esctl get snapshot-repo <repo>
 esctl get snapshot-repos
 esctl set snapshot <repo> <snapshot>
 esctl update snapshot <repo> <snapshot>
 esctl delete snapshot <repo> <snapshot>
 esctl set snapshot-repo <repo> --type=fs --settings="location:/backup"
 esctl delete snapshot-repo <repo>
 
 # Security
 esctl get users
 esctl get user <username>
 esctl set user <username> --file=user.json
 esctl delete user <username>
 esctl get roles
 esctl get role <role>
 esctl set role <role> --file=role.json
 esctl delete role <role>
 
 # Reindex
 esctl set reindex --source=<index> --dest=<index>
 esctl get reindex <task-id>
 esctl delete reindex <task-id>
 ```
