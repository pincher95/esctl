package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAliasesCmd = &cobra.Command{
	Use:                   "aliases [--index index]",
	DisableFlagsInUseLine: true,
	Short:                 "Retrieves information for one or more data stream or index aliases.",
	Long: utils.Trim(`
	Get Elasticsearch aliases. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all aliases.
	esctl get aliases

	# Retrieve aliases for a specific index.
	esctl get aliases --index my_index
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		aliasClient := index.NewIndex()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		if flagIndex != "" {
			if err := validation.ValidateIndexPattern(flagIndex); err != nil {
				return err
			}
		}

		if !flagRefresh {
			return handleAliasLogic(ctx, aliasClient, *conf)
		}

		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleAliasLogic(ctx, aliasClient, *conf)
		})
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
	{Header: "ROUTING", Type: output.Text},
	{Header: "INDEX-ROUTING", Type: output.Text},
	{Header: "SEARCH-ROUTING", Type: output.Text},
	{Header: "IS_WRITE_INDEX", Type: output.Boolean},
}

func handleAliasLogic(ctx context.Context, client index.Index, conf config.Config) error {
	aliases, err := client.GetAliases(ctx, flagIndex)
	if err != nil {
		return fmt.Errorf("failed to retrieve aliases: %w", err)
	}

	columnDefs, err := getColumnDefs(conf, "alias", aliasColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := [][]string{}

	for idx, detail := range *aliases {
		for alias, aliasDetails := range detail.Aliases {
			if includeIndexByWriteIndex(aliasDetails) {
				rowData := map[string]string{
					"ALIAS":          alias,
					"INDEX":          idx,
					"INDEX-ROUTING":  aliasDetails.IndexRouting,
					"SEARCH-ROUTING": aliasDetails.SearchRouting,
					"ROUTING":        aliasDetails.Routing,
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

	return output.PrintTable(columnDefs, data, sortCols)
}

func includeIndexByWriteIndex(aliasDetails index.AliasDetails) bool {
	switch {
	case flagWritable && aliasDetails.IsWriteIndex:
		return true
	case !flagWritable:
		return true
	}
	return false
}
