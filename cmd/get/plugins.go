package get

import (
	"context"
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginsClient := cat.NewCat()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		if !flagRefresh {
			return handlePluginsLogic(ctx, pluginsClient, *conf)
		}

		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handlePluginsLogic(ctx, pluginsClient, *conf)
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

func handlePluginsLogic(ctx context.Context, client cat.Cat, conf config.Config) error {
	plugins, err := client.CatPlugins(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to retrieve plugins: %w", err)
	}

	columnDefs, err := getColumnDefs(conf, "plugins", pluginsColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
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
		return output.PrintTable(columnDefs, data, sortCols)
	}
	sortCols := output.ParseSortColumns("NAME")
	return output.PrintTable(columnDefs, data, sortCols)
}
