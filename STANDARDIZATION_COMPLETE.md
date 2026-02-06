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
esctl get alias --name <alias>
esctl set alias --name <alias> --indices=<index1,index2,...>
esctl delete alias --name <alias> --indices=<index1,index2,...>
esctl update alias --name <alias> --from=<from-index> --to=<to-index>
 
 # Pipelines
 esctl get pipelines
esctl get pipeline --id <pipeline-id>
esctl set pipeline --id <pipeline-id> --file=pipeline.json
esctl delete pipeline --id <pipeline-id>
 esctl update pipeline --file=request.json [--pipeline=<pipeline-id>]
 
 # Snapshots
esctl get snapshot --repository <repo> --name <snapshot>
esctl get snapshot --repository <repo>
esctl get snapshot-status [--repository <repo>] [--name <snapshot>]
esctl get snapshot-repo --repository <repo>
 esctl get snapshot-repos
esctl set snapshot --repository <repo> --name <snapshot>
esctl update snapshot --repository <repo> --name <snapshot>
esctl delete snapshot --repository <repo> --name <snapshot>
esctl set snapshot-repo --repository <repo> --type=fs --settings="location:/backup"
esctl delete snapshot-repo --repository <repo>
 
 # Security
 esctl get users
esctl get user --name <username>
esctl set user --name <username> --file=user.json
esctl delete user --name <username>
 esctl get roles
esctl get role --name <role>
esctl set role --name <role> --file=role.json
esctl delete role --name <role>
 
 # Reindex
 esctl set reindex --source=<index> --dest=<index>
esctl get reindex --task-id <task-id>
esctl delete reindex --task-id <task-id>
 ```
