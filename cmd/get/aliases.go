package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAliasesCmd = &cobra.Command{
	Use:                   "aliases",
	DisableFlagsInUseLine: true,
	Short:                 "Retrieves information for one or more data stream or index aliases.",
	Args:                  cobra.NoArgs,
	Long: utils.Trim(`
	Get Elasticsearch aliases. You can filter the results using the index or name flags.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all aliases.
	esctl get aliases

	# Retrieve a specific alias by name.
	esctl get aliases --name my-alias

	# Retrieve aliases for a specific index.
	esctl get aliases --index my_index

	# Retrieve aliases by name substring.
	esctl get aliases --name logs
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// If a specific alias name is provided, get that alias
		if flagAliasName != "" {
			return handleGetSpecificAlias(ctx, flagAliasName)
		}

		// Otherwise, list all aliases
		aliasClient := index.NewIndex()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

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
	getAliasesCmd.Flags().StringVar(&flagAliasName, "name", "", "Alias name to retrieve or substring for filtering")
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

	// Calculate total capacity for pre-allocation
	totalAliases := 0
	for _, detail := range *aliases {
		totalAliases += len(detail.Aliases)
	}
	data := make([][]string, 0, totalAliases)

	for idx, detail := range *aliases {
		for alias, aliasDetails := range detail.Aliases {
			if flagAliasName != "" && !strings.Contains(alias, flagAliasName) {
				continue
			}
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

func handleGetSpecificAlias(ctx context.Context, aliasName string) error {
	if err := validation.ValidateAliasName(aliasName); err != nil {
		return err
	}

	aliases, err := index.GetAlias(ctx, nil, []string{aliasName})
	if err != nil {
		return fmt.Errorf("failed to get alias: %w", err)
	}

	if len(aliases) == 0 {
		return fmt.Errorf("alias not found: %s", aliasName)
	}

	return output.Render(aliases)
}

var flagAliasName string
