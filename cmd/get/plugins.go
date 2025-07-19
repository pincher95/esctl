package get

import (
	"context"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	cat "github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getPluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Get plugin information for the Elasticsearch cluster",
	Long: utils.Trim(`
	Get a list of plugins running on each node of a cluster.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all plugins.
	esctl get plugins
	`),
	Run: func(cmd *cobra.Command, args []string) {
		pluginsClient := cat.NewCat()
		conf := config.ParseConfigFile()

		ctx := cmd.Context()

		if !flagRefresh {
			handlePluginsLogic(ctx, pluginsClient, *conf)
			return
		}

		utils.WatchLoop(flagRefreshInterval, func() error {
			handlePluginsLogic(ctx, pluginsClient, *conf)
			return nil
		})
	},
}

func init() {
}

var pluginsColumns = []output.ColumnDefaults{
	{Header: "ID", Type: output.Text},
	{Header: "NAME", Type: output.Text},
	{Header: "COMPONENT", Type: output.Text},
	{Header: "VERSION", Type: output.Text},
	{Header: "DESCRIPTION", Type: output.Text},
}

func handlePluginsLogic(ctx context.Context, client cat.Cat, conf config.Config) {
	plugins, err := client.CatPlugins(ctx, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve plugins:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "plugins", pluginsColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, plugin := range plugins {
		rowData := map[string]string{
			"NAME":        plugin.Name,
			"COMPONENT":   plugin.Component,
			"VERSION":     plugin.Version,
			"DESCRIPTION": plugin.Description,
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
