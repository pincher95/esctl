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
esctl get alias --name <alias>
esctl set alias --name <alias> --indices=<index1,index2,...>
esctl delete alias --name <alias> --indices=<index1,index2,...>
esctl update alias --name <alias> --from=<from-index> --to=<to-index>
 ```

 ### Pipelines
 ```bash
 esctl get pipelines
esctl get pipeline --id <pipeline-id>
esctl set pipeline --id <pipeline-id> --file=pipeline.json
esctl delete pipeline --id <pipeline-id>
 esctl update pipeline --file=request.json [--pipeline=<pipeline-id>]
 ```

 ### Snapshots
 ```bash
esctl get snapshot --repository <repo> --name <snapshot>       # Get details
esctl get snapshot --repository <repo>                         # List snapshots in repo
esctl get snapshot-status [--repository <repo>] [--name <snapshot>]
esctl get snapshot-repos
esctl get snapshot-repos --name <repo>          # Filter repositories

esctl set snapshot --repository <repo> --name <snapshot>
esctl update snapshot --repository <repo> --name <snapshot>    # Restore
esctl delete snapshot --repository <repo> --name <snapshot>

esctl set snapshot-repo --repository <repo> --type=fs --settings="location:/backup"
esctl delete snapshot-repo --repository <repo>
 ```

 ### Security (users/roles)
 ```bash
 esctl get users
esctl get user --name <username>
esctl set user --name <username> --file=user.json
esctl delete user --name <username>

 esctl get roles
esctl get role --name <role>
esctl set role --name <role> --file=role.json
esctl delete role --name <role>
 ```

 ### Reindex
 ```bash
 esctl set reindex --source=<index> --dest=<index>
esctl get reindex --task-id <task-id>
esctl delete reindex --task-id <task-id>
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
