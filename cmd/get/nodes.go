package get

import (
	"fmt"
	"os"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Get all nodes in the Elasticsearch cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := cmd.Context()
		config := config.ParseConfigFile()

		// select
		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleNodeLogic(*config)
			return nil
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleNodeLogic(*config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getNodesCmd.Flags().StringVar(&flagNode, "node", "", "Filter shards by node")
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

func handleNodeLogic(conf config.Config) {
	nodes, err := es.GetNodes(flagNode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve nodes: %v\n", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "node", nodeColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, node := range nodes {
		rowData := map[string]string{
			"IP":           node.IP,
			"HEAP-PERCENT": node.HeapPercent + "%",
			"RAM-PERCENT":  node.RAMPercent + "%",
			"CPU":          node.CPU + "%",
			"LOAD-1M":      node.Load1m,
			"LOAD-5M":      node.Load5m,
			"LOAD-15M":     node.Load15m,
			"NODE-ROLE":    node.NodeRole,
			"NODE-ROLES":   node.NodeRoles,
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
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("NAME")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
