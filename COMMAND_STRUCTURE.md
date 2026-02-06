# esctl Command Structure
 
 This document reflects the current `esctl` command structure.
 
 ## Core verbs
 
 `esctl` uses verb-first commands for resource operations:
 
 - `get` — read or list resources
 - `set` — create or update resources
 - `delete` — remove resources
 - `update` — operational updates (reroute, cache clear, alias move, snapshot restore, pipeline simulate, index settings)
 
 ## Common resources
 
 ### Aliases
 ```bash
 esctl get aliases
 esctl get alias <alias>
 esctl set alias <alias> --indices=<index1,index2,...>
 esctl delete alias <alias> --indices=<index1,index2,...>
 esctl update alias <alias> --from=<from-index> --to=<to-index>
 ```
 
 ### Pipelines
 ```bash
 esctl get pipelines
 esctl get pipeline <pipeline-id>
 esctl set pipeline <pipeline-id> --file=pipeline.json
 esctl delete pipeline <pipeline-id>
 esctl update pipeline --file=request.json [--pipeline=<pipeline-id>]
 ```
 
 ### Snapshots
 ```bash
 esctl get snapshot <repo> <snapshot>       # Get details
 esctl get snapshot <repo>                  # List snapshots in repo
 esctl get snapshot --repository <repo>     # List snapshots in repo
 esctl get snapshot-status [repo] [snapshot]
 esctl get snapshot-repo <repo>
 esctl get snapshot-repos
 
 esctl set snapshot <repo> <snapshot>
 esctl update snapshot <repo> <snapshot>    # Restore
 esctl delete snapshot <repo> <snapshot>
 
 esctl set snapshot-repo <repo> --type=fs --settings="location:/backup"
 esctl delete snapshot-repo <repo>
 ```
 
 ### Security (users/roles)
 ```bash
 esctl get users
 esctl get user <username>
 esctl set user <username> --file=user.json
 esctl delete user <username>
 
 esctl get roles
 esctl get role <role>
 esctl set role <role> --file=role.json
 esctl delete role <role>
 ```
 
 ### Reindex
 ```bash
 esctl set reindex --source=<index> --dest=<index>
 esctl get reindex <task-id>
 esctl delete reindex <task-id>
 ```
 
 ## Other top-level commands
 
 These remain as dedicated commands:
 
 - `template <action>` — template management
 - `ilm <action>` — ILM policy management
 - `index <action>` — index maintenance (refresh/flush/forcemerge)
 - `describe <entity>` — detailed JSON/YAML views
 - `count`, `query`, `bulk`, `analyze`, `explain`, `profile`, `tasks`, `version`, `config`
 
 ## Global flags
 
 ```bash
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
