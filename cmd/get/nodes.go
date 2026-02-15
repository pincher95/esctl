package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var getNodesCmd = &cobra.Command{
	Use:                   "nodes [--name node] [--bytes unit] [--time unit]",
	DisableFlagsInUseLine: true,
	Short:                 "Get information about the nodes for the Elasticsearch cluster.",
	Long: utils.Trim(`
Get information about the nodes in a cluster. You can filter the results using the name flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all nodes.
	esctl get nodes

	# Retrieve nodes for a specific node.
	esctl get nodes --name my_node

	# Retrieve nodes in kilobytes.
	esctl get nodes --bytes kb

	# Retrieve nodes time in seconds.
	esctl get nodes --bytes kb --time s
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeClient := cat.NewCat()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		if !flagRefresh {
			return handleNodeLogic(ctx, nodeClient, *conf)
		}

		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleNodeLogic(ctx, nodeClient, *conf)
		})
	},
}

func init() {
	getNodesCmd.Flags().StringVar(&flagFilter, "name", "", "Filter by node name or substring of node name e.g. 'data-node-1', 'node', 'data'")
	getNodesCmd.Flags().StringVar(&flagBytes, "bytes", "", "The unit in which to display byte values. Valid values are: 'b', 'kb', 'mb', 'gb', 'tb', 'pb'.")
	getNodesCmd.Flags().StringVar(&flagTime, "time", "", "Specifies the time units, for example, 5d or 7h. Valid values are: nanos, micros, ms, s, m, h, d.")
}

var nodeColumns = []output.ColumnDefaults{
	{Header: "IP", Type: output.Text},
	{Header: "HEAP-PERCENT", Type: output.Percent},
	{Header: "RAM-PERCENT", Type: output.Percent},
	{Header: "CPU", Type: output.Percent},
	{Header: "LOAD-1M", Type: output.Number},
	{Header: "LOAD-5M", Type: output.Number},
	{Header: "LOAD-15M", Type: output.Number},
	{Header: "NODE-ROLE", Type: output.Text},
	{Header: "NODE-ROLES", Type: output.Text},
	{Header: "MASTER", Type: output.Text},
	{Header: "NAME", Type: output.Text},
}

func handleNodeLogic(ctx context.Context, client cat.Cat, conf config.Config) error {
	nodes, err := client.CatNodes(ctx, "", flagFilter, flagBytes, flagTime)
	if err != nil {
		return fmt.Errorf("failed to retrieve nodes: %w", err)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(nodes)
	}

	columnDefs, err := getColumnDefs(conf, "node", nodeColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := make([][]string, 0, len(nodes))

	for _, node := range nodes {
		rowData := map[string]string{
			"IP":           node.IP,
			"HEAP-PERCENT": fmt.Sprintf("%d%%", node.HeapPercent),
			"RAM-PERCENT":  fmt.Sprintf("%d%%", node.RAMPercent),
			"CPU":          fmt.Sprintf("%d%%", node.CPU),
			"LOAD-1M":      node.Load1M,
			"LOAD-5M":      node.Load5M,
			"LOAD-15M":     node.Load15M,
			"NODE-ROLE":    node.Role,
			"NODE-ROLES":   node.Roles,
			"MASTER":       node.Master,
			"NAME":         node.Name,
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		return output.PrintTable(columnDefs, data, sortCols)
	}
	sortCols := output.ParseSortColumns("NAME")
	return output.PrintTable(columnDefs, data, sortCols)
}
