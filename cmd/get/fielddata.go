package get

import (
	"fmt"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getFielddataCmd = &cobra.Command{
	Use:   "fielddata",
	Short: "Get Elasticsearch fielddata",
	Long: utils.Trim(`
	Get the amount of heap memory currently used by the field data cache on every data node in the cluster.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all fielddata.
	esctl get fielddata

	# Retrieve fielddata for a specific fields.
	esctl get fielddata --fields my_field1,my_field2

	# Retrieve fielddata for a specific fields and display the size in kilobytes.
	esctl get fielddata --fields my_field1,my_field2 --bytes kb
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		feilddataClient := cat.NewCat()
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			return handleFielddataLogic(feilddataClient, *config)
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			err := handleFielddataLogic(feilddataClient, *config)
			if err != nil {
				return err
			}
			time.Sleep(flagRefreshInterval)
		}
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

func handleFielddataLogic(client cat.Cat, conf config.Config) error {
	fielddata, err := client.CatFielddata(nil, &flagFields, &flagBytes)
	if err != nil {
		return fmt.Errorf("Failed to retrieve fielddata: %v", err)
	}

	columnDefs, err := getColumnDefs(conf, "id", fielddataColumns)
	if err != nil {
		return fmt.Errorf("Failed to get column definitions: %v", err)
	}

	data := [][]string{}

	if fielddata != nil {
		for _, feild := range *fielddata {
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
	}

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("id")
		output.PrintTable(columnDefs, data, sortCols)
	}

	return nil
}
