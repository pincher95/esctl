package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getFielddataCmd = &cobra.Command{
	Use:                   "fielddata [--fields fields] [--bytes bytes]",
	DisableFlagsInUseLine: true,
	Short:                 "Get field data cache information for the Elasticsearch cluster",
	Long: utils.Trim(`
	Get the amount of heap memory currently used by the field data cache on every data node in the cluster. The field data cache is used to cache the field values of fields used in filters and facets. This command provides information about the field data cache size for each field on each node in the cluster.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all fielddata.
	esctl get fielddata

	# Retrieve fielddata for a specific fields.
	esctl get fielddata --fields my_field1,my_field2

	# Retrieve fielddata display the size in kilobytes.
	esctl get fielddata --bytes kb
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		fielddataClient := cat.NewCat()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		if !flagRefresh {
			return handleFielddataLogic(ctx, fielddataClient, *conf)
		}

		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleFielddataLogic(ctx, fielddataClient, *conf)
		})
	},
}

func init() {
	getFielddataCmd.Flags().StringVar(&flagFields, "fields", "", "Comma-separated list of fields used to limit returned information. To retrieve all fields, omit this parameter.")
	getFielddataCmd.Flags().StringVar(&flagBytes, "bytes", "", "The unit in which to display byte values. Valid values are: 'b', 'kb', 'mb', 'gb', 'tb', 'pb'.")
}

var fielddataColumns = []output.ColumnDefaults{
	{Header: "ID", Type: output.Text},
	{Header: "HOST", Type: output.Text},
	{Header: "IP", Type: output.Text},
	{Header: "NODE", Type: output.Text},
	{Header: "FIELD", Type: output.Text},
	{Header: "SIZE", Type: output.Text},
}

func handleFielddataLogic(ctx context.Context, client cat.Cat, conf config.Config) error {
	fielddata, err := client.CatFielddata(ctx, "", flagFields, flagBytes)
	if err != nil {
		return fmt.Errorf("failed to retrieve fielddata: %w", err)
	}

	columnDefs, err := getColumnDefs(conf, "id", fielddataColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := [][]string{}

	for _, feild := range fielddata {
		rowData := map[string]string{
			// "ID":    feild.ID,
			"HOST":  feild.Host,
			"IP":    feild.IP,
			"NODE":  feild.Node,
			"FIELD": feild.Field,
			"SIZE":  feild.Size,
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
	sortCols := output.ParseSortColumns("id")
	return output.PrintTable(columnDefs, data, sortCols)
}
