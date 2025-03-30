package get

import (
	"fmt"
	"os"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAliasesCmd = &cobra.Command{
	Use:                   "aliases [--index index]",
	DisableFlagsInUseLine: true,
	Short:                 "Get Elasticsearch aliases for the Elasticsearch cluster",
	Long: utils.Trim(`
	Get Elasticsearch aliases. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all aliases.
	esctl get aliases

	# Retrieve aliases for a specific index.
	esctl get aliases --index my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		aliasClient := index.NewIndex()

		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleAliasLogic(aliasClient, *config)
			return
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleAliasLogic(aliasClient, *config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getAliasesCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
	getAliasesCmd.Flags().BoolVar(&flagWritable, "writable", true, "Filter by writable index")
}

var aliasColumns = []output.ColumnDefaults{
	{Header: "ALIAS", Type: output.Text},
	{Header: "INDEX", Type: output.Text},
	{Header: "FILTER", Type: output.Text},
	{Header: "ROUTING-INDEX", Type: output.Text},
	{Header: "ROUTING-SEARCH", Type: output.Text},
	{Header: "IS_WRITE_INDEX", Type: output.Boolean},
}

func handleAliasLogic(client index.Index, conf config.Config) {
	aliases, err := client.GetAliases(nil, &flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve aliases:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "alias", aliasColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for index, detail := range *aliases {
		for alias := range detail.Aliases {
			if includeIndexByWriteIndex(detail.Aliases[alias]) {
				rowData := map[string]string{
					"ALIAS": alias,
					"INDEX": index,
					// "FILTER":         detail.Aliases[alias].Filter,
					"ROUTING-INDEX":  detail.Aliases[alias].IndexRouting,
					"ROUTING-SEARCH": detail.Aliases[alias].SearchRouting,
					// "IS_WRITE_INDEX": detail.Aliases[alias].IsWriteIndex,
				}

				row := make([]string, len(columnDefs))
				for i, colDef := range columnDefs {
					row[i] = rowData[colDef.Header]
				}
				data = append(data, row)
			}
		}
	}

	sortCols := output.ParseSortColumns(flagSortBy)

	output.PrintTable(columnDefs, data, sortCols)
}

func includeIndexByWriteIndex(aliasInfo index.AliasInfo) bool {
	switch {
	case flagWritable && aliasInfo.IsWriteIndex:
		return true
	case !flagWritable:
		return true
	}
	return false
}
