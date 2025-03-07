# esctl

[![Latest Stable Version](https://img.shields.io/badge/version-v1.0.0-blue.svg)](https://github.com/pincher95/esctl/releases/tag/v1.0.0)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/pincher95/esctl/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pincher95/esctl)](https://goreportcard.com/report/github.com/pincher95/esctl)

`esctl` is a command-line tool for managing and interacting with Elasticsearch clusters. It provides a convenient interface to perform various operations, such as querying cluster information, managing indices, retrieving shard details, and monitoring tasks.

## Features:
- Retrieve information about nodes, indices, shards, aliases, and tasks in an Elasticsearch cluster
- Describe cluster health and stats
- Describe index settings and mappings
- Simple and intuitive command-line interface

## Contributing
Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## Table of Contents

- [Installation](#installation)
- [Examples](#examples)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Get](#get)
  - [Describe](#describe)
  - [Count](#count)
  - [Count with Grouping](#count-with-grouping)
  - [Query](#query)
- [License](#license)

## Installation

To install `esctl`, ensure that you have Go installed and set up in your development environment. Then, follow the steps below:

1. Open a terminal or command prompt.

2. Run the following command to install `esctl`:

   ```shell
   go install github.com/pincher95/esctl@latest
   ```
   This command will fetch the source code from the GitHub repository, compile it, and install the `esctl` binary in your Go workspace's `bin` directory.

3. Make sure that your Go workspace's `bin` directory is added to your system's `PATH` environment variable. This step will allow you to run `esctl` from any directory in the terminal or command prompt.

Once installed, you can run `esctl` by simply typing `esctl` in the terminal or command prompt.

## Examples

```shell
> esctl get shards --index=articles --primary
INDEX     ID                      SHARD  PRI-REP  STATE    DOCS  STORE  IP         NODE               SEGMENTS-COUNT
articles  jxn-Oa3XSPigaCBYt9fKiw  0      primary  STARTED  0     225b   127.0.0.1  es-data-0          0
articles  jxn-Oa3XSPigaCBYt9fKiw  1      primary  STARTED  0     225b   127.0.0.1  es-data-0          0
articles  jxn-Oa3XSPigaCBYt9fKiw  2      primary  STARTED  0     225b   127.0.0.1  es-data-0          0
```

```shell
> esctl get shards --index=articles --shard 0 --unassigned --sort-by=unassigned-at
INDEX     SHARD  PRI-REP  STATE       UNASSIGNED-REASON  UNASSIGNED-AT
articles  0      replica  UNASSIGNED  CLUSTER_RECOVERED  2023-05-07T20:37:07.520Z
articles  0      replica  UNASSIGNED  CLUSTER_RECOVERED  2023-05-07T20:37:07.520Z
```

```shell
> esctl get indices
HEALTH  STATUS  INDEX     UUID                    PRI  REP  DOCS-COUNT  DOCS-DELETED  CREATION-DATE             STORE-SIZE  PRI-STORE-SIZE
yellow  open    articles  8vCars4rQquYHNhpKV2fow  3    2    0           0             2023-05-07T19:17:52.259Z  675b        675b
```

```shell
> esctl get nodes --sort-by=cpu
NAME               IP         NODE-ROLE    MASTER  HEAP-MAX  HEAP-CURRENT  HEAP-PERCENT  CPU  LOAD-1M  DISK-TOTAL  DISK-USED  DISK-AVAILABLE
es-data-0          127.0.0.1  cdfhilmrstw  *       4gb       1.6gb         41%           10%  2.02     232.9gb     199.2gb    33.6gb
```

```shell
> esctl get aliases --index=articles
ALIAS           INDEX
articles_alias  articles
```

## Configuration

`esctl` supports reading context configurations from a YAML file, enabling you to easily switch between different Elasticsearch contexts.

Create a configuration file named `esctl.yml` in your `$HOME/.config` directory. This file will contain the details of the Elasticsearch contexts you want to connect to.

Here is an example configuration:

```yaml
current-context: "local"
contexts:
  - name: "local"
    protocol: "http"
    host: "localhost"
    port: 9200
  - name: "production"
    protocol: "https"
    host: "prod.es.example.com"
    port: 443
    username: "prod_username"
    password: "prod_password"
```

In the configuration file:

- `current-context` is the name of the context that `esctl` will connect to by default.
- `contexts` is an array of your Elasticsearch contexts.
  - `name` is the name you assign to the context.
  - `protocol`, `host`, `port`, `username`, and `password` are the connection details for each context.
  - `protocol` and `port` are optional and default to `http` and `9200` respectively.

> **Note**<br>
> `esctl` will use the `current-context` defined in the configuration file unless another cluster is specified via command-line flag or environment variable.

### esctl config use-context

Sets the current context in the configuration file.

```bash
esctl config use-context CONTEXT
```

- `CONTEXT`: The name of the context to set as the current context.

This command updates the `current-context` in the configuration file (`esctl.yml`) with the specified context name. The updated configuration will be used for subsequent operations performed by `esctl`.

> **Note**<br>
> The specified context name must already be defined in the configuration file.

### esctl config get-contexts

Displays the contexts defined in the `esctl.yml` file.

```bash
esctl config get-contexts
```

Example output:

```yaml
- name: local(*)
  protocol: https
  host: localhost
  port: 9200
  username: myuser
  password: ********
- name: production
  protocol: http
  host: example.com
  port: 9200
  username: anotheruser
  password: ********
```

In the example output, the current context is marked with an asterisk (*).

### esctl config current-context

Displays the current context that `esctl` is configured to use.

```bash
esctl config current-context
```

Example output:

```
local
```

### Displaying Current Context with Starship Prompt

If you're using [Starship](https://starship.rs) for your prompt, you can display the current `esctl` context directly in your shell prompt. This can be done by adding a custom module to your Starship configuration. Here's an example:

```toml
[custom.esctl]
command = "esctl config current-context 2>/dev/null || echo 'none'"
description = "Displays the current esctl context"
when = "command -v esctl"
symbol = "🄴 "
```

This configuration adds a new custom module that executes `esctl config current-context` and displays the current context in your prompt. The module will only appear if `esctl` is installed (`when = "command -v esctl"`). The `symbol` option is used to provide a visual indicator for the esctl context information. The `🄴` symbol is used to represent Elasticsearch.


### Elasticsearch Host Configuration

Additionally, `esctl` allows you to configure the Elasticsearch host, port, protocol, username, and password using command-line flags or environment variables. By default, the port is set to `9200`, and the protocol to `http`.

To specify a custom host, you can use the `--host` flag followed by the desired host value. For example:

```shell
esctl --host=HOST COMMAND
```

Similarly, to specify a custom port, you can use the `--port` flag followed by the desired port value. For example:

```shell
esctl --port=PORT COMMAND
```
To specify a custom protocol, you can use the `--protocol` flag followed by either `http` or `https`. For example:

```shell
esctl --protocol=https COMMAND
```

To provide basic authentication credentials, you can use the `--username` and `--password` flags followed by the corresponding values. For example:

```shell
esctl --username=USERNAME --password=PASSWORD COMMAND
```

Alternatively, you can set the `ESCTL_HOST`, `ESCTL_PORT`, `ESCTL_PROTOCOL`, `ESCTL_USERNAME`, and `ESCTL_PASSWORD` environment variables to your desired Elasticsearch configuration.

If the corresponding command-line flags and environment variables are not provided, `esctl` will use the default values (`9200`, `http`, no username, and no password) for the Elasticsearch connection.

> **Warning**<br>
> Since host is mandatory, if host is not provided via a flag, environment variable or esctl.yml, esctl will exit with an error.

### Customizing Columns

You can customize the columns displayed when running `esctl get ENTITY` using the `esctl.yml` configuration file.

To customize the columns, add an optional `entities` field to the `esctl.yml` file. Under `entities`, specify the desired entities (`node`, `index`, `shard`, `alias`, `task`) and their corresponding columns. Here is an example:

```yaml
contexts:
  - host: localhost
    name: cluster1
    port: 9200
    protocol: http
  - host: 127.0.0.1
    name: cluster2
current-context: cluster2
entities:
  node:
    columns:
      - "NAME"
      - "IP"
      - "NODE-ROLE"
      - "MASTER"
      - "HEAP-MAX"
      - "HEAP-CURRENT"
      - "CPU"
      - "LOAD-1M"
  index:
    columns: []
  task:
    columns: []
```

In the `columns` field of each entity, specify the desired columns in the order you want them to appear.

> **Note**<br>
> If you do not provide a `columns` field for an entity, it will use the default columns.

## Usage

### Get

The `get` command allows you to retrieve information about Elasticsearch entities. Supported entities include nodes, indices, shards, aliases, and tasks. This command provides a read-only view of the cluster and does not support data querying.

```shell
esctl get ENTITY [flags]
```

#### Available Entities

- `nodes`: List all nodes in the Elasticsearch cluster.
- `indices`: List all indices in the Elasticsearch cluster.
- `shards`: List detailed information about shards, including their sizes and placement.
- `aliases`: List all aliases in the Elasticsearch cluster.
- `tasks`: List all tasks in the Elasticsearch cluster.

#### Flags

- `--index`: Specifies the name of the index (applies to `indices`, `shards`, and `aliases` entities).
- `--node`: Specified the name of the node (applies to `nodes` and `shards` entities).
- `--shard`: Filters shards by shard number.
- `--primary`: Filters primary shards.
- `--replica`: Filters replica shards.
- `--started`: Filters shards in STARTED state.
- `--relocating`: Filters shards in RELOCATING state.
- `--initializing`: Filters shards in INITIALIZING state.
- `--unassigned`: Filters shards in UNASSIGNED state.
- `--actions`: Filters tasks by actions.
- `--sort-by`: Specifies the columns to sort by, separated by commas (applies to all entities). The column names are case insensitive.
- `--columns`: Specifies the columns to display, separated by commas (applies to all entities). To display all columns, use `all`. The column names are case insensitive.

#### Get Nodes

Retrieves a list of all nodes in the Elasticsearch cluster.

```shell
esctl get nodes
```

#### Get Indices

Retrieves a list of all indices in the Elasticsearch cluster.

```shell
esctl get indices
```

#### Get Shards

To retrieve shards from Elasticsearch, you can use the following command:

```shell
esctl get shards [--index index] [--node node] [--shard shard] [--primary] [--replica] [--started] [--relocating] [--initializing] [--unassigned]
```

* `--index`: Specifies the name of the index to retrieve shards from.
* `--node`: Filters shards by node name.
* `--shard`: Filters shards by shard number.
* `--primary`: Filters primary shards.
* `--replica`: Filters replica shards.
* `--started`: Filters shards in the STARTED state.
* `--relocating`: Filters shards in the RELOCATING state.
* `--initializing`: Filters shards in the INITIALIZING state.
* `--unassigned`: Filters shards in the UNASSIGNED state.

If none of the flags are provided, all shards will be returned.

Example usage:

```shell
esctl get shards --index my_index --relocating
```
This will retrieve only the shards that are currently relocating for the specified index.

#### Get Aliases

Retrieves the list of aliases defined in Elasticsearch, including the index names they are associated with.

Usage:

```shell
esctl get aliases [--index INDEX]
```

`--index`: Filter the aliases by a specific index. If not provided, aliases to all indices will be returned.

#### Get Tasks

The `get tasks` command retrieves information about tasks in the Elasticsearch cluster.

Usage:

```shell
esctl get tasks [--actions ACTIONS]
```

Example:

```shell
esctl get tasks --actions 'index*' --actions '*search*'

```

### Describe

The `esctl describe` command allows you to retrieve detailed information about various entities in the Elasticsearch cluster. The output is in JSON or YAML format, making it easy to read and understand. You can select your preferred output format using the `--output` or `-o` flag, with `json` and `yaml` being the available options.

#### Describe Cluster

This command outputs the cluster information in YAML format, providing a comprehensive overview of the cluster's current state.

```shell
esctl describe cluster
```

#### Describe Index

This command outputs the mappings and settings of a specified index in JSON format.

```shell
esctl describe index INDEX
```

This command also supports the `--mappings` and `--settings` flags, which can be used to get only the mappings or settings respectively.

```shell
# To get only mappings
esctl describe index INDEX --mappings

# To get only settings
esctl describe index INDEX --settings
```

> **Note**<br>
> Consider piping the output of `describe index` to [fx](https://github.com/antonmedv/fx), a command-line JSON processing tool, for a more convenient experience.

```shell
esctl describe index INDEX | fx
```

### Count

![esctl usage](./assets/count.gif)

The `esctl count` command allows you to retrieve the count of documents in one or more Elasticsearch indices. You can specify filters to apply using the `--term` and `--exists` flags.

#### Count All Documents

To count all documents across all indices, use the following command:

```shell
esctl count
```

#### Count Documents in Specific Index

To count all documents in a specific index, use the following command:

```shell
esctl count --index index
```

#### Count Documents with Term Filters

You can apply term filters to count documents that match specific field-value combinations. Use the `--term` flag followed by the field-value pairs separated by a colon (`:`). For example, to count documents with the field `price` equal to `12` and `category` equal to `electronics`, use the following command:

```shell
esctl count --term "price:12" --term "category:electronics"
```

#### Count Documents with Existence Filters

To count documents based on the existence of a field, use the `--exists` flag followed by the field name. For example, to count documents where the field `category` exists, use the following command:

```shell
esctl count --exists "category"
```

> **Note**<br>
> You can combine both term and existence filters in a single command to further refine the count.

### Count with Grouping

The `esctl count` command also supports grouping the documents by a specific field and displaying the respective counts. You can use the `--group-by` flag to specify the field to group by.

#### Group Documents by Field

To count and group documents by a specific field, use the `--group-by` flag followed by the field name. For example, to count documents in the index `articles` and group them by the `category` field, use the following command:

```shell
esctl count --index articles --group-by category
```

This command will retrieve the count of documents in the `articles` index and group them based on the values of the `category` field.

#### Grouping with Filters

You can combine the grouping functionality with term and existence filters to further refine the count and group the documents accordingly. For example, to count and group documents in the index `articles` with the field `price` equal to `12`, use the following command:

```shell
esctl count --index articles --term "price:12" --group-by category
```

This command will count the documents in the `articles` index that satisfy the term filter (`price:12`) and group them by the values of the `category` field.

> **Note**<br>
> If an index name is not provided to the count command with `--group-by`, all the indices will be grouped individually.

### Query

The `query` command allows you to execute queries against Elasticsearch.

```sh
esctl query INDEX
```

#### Flags

- `--id`: Specify document IDs to fetch. Can be specified multiple times.

  Example: `--id 61 --id 62`

- `--term (-t)`: Term filters to apply. The format should be `field:value`. Can be specified multiple times.

  Example: `--term "price:10" --term "category:electronics"`

- `--size`: Specify the number of hits to return. Defaults to 1.

  Example: `--size 5`

#### Examples

```sh
esctl query articles
esctl query articles --id 61
esctl query articles --term "price:10" --size 2
```

This would respectively:

- Query all documents in the `articles` index.
- Query the `articles` index and get the document with ID `61`.
- Query the `articles` index filtering by the term `price:10` and return 2 hits.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
