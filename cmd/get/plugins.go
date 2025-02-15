package get

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getPluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Get Elasticsearch plugins",
	Long: utils.Trim(`
	Get Elasticsearch plugins.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all plugins.
	esctl get plugins
	`),
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.ParseConfigFile()
		handlePluginsLogic(*conf)
	},
}

func init() {
}

var pluginsColumns = []output.ColumnDef{
	{Header: "ID", Type: output.Text},
	{Header: "NAME", Type: output.Text},
	{Header: "COMPONENT", Type: output.Text},
	{Header: "VERSION", Type: output.Text},
	{Header: "DESCRIPTION", Type: output.Text},
}

func handlePluginsLogic(conf config.Config) {
	plugins, err := es.GetPlugins()
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
			"ID":          plugin.ID,
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
		output.PrintTable(columnDefs, data, flagSortBy...)
	} else {
		output.PrintTable(columnDefs, data, "NAME")
	}
}
