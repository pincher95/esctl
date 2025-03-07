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

var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Get tasks information",
	Long:  `This command retrieves and displays tasks information from Elasticsearch cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleTaskLogic(*config)
			return
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleTaskLogic(*config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getTasksCmd.Flags().StringArrayVarP(&flagActions, "actions", "a", nil, "Filter tasks by actions")
}

var taskColumns = []output.ColumnDefaults{
	{Header: "NODE", Type: output.Text},
	{Header: "ID", Type: output.Number},
	{Header: "ACTION", Type: output.Text},
	{Header: "DESCRIPTION", Type: output.Text},
	{Header: "START-TIME", Type: output.Number},
	{Header: "RUNNING-TIME", Type: output.Number},
}

func handleTaskLogic(config config.Config) {
	tasksResponse, err := es.GetTasks(flagActions)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve tasks:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(config, "task", taskColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, node := range tasksResponse.Nodes {
		for _, task := range node.Tasks {
			rowData := map[string]string{
				"NODE":         task.Node,
				"ID":           fmt.Sprintf("%d", task.ID),
				"ACTION":       task.Action,
				"DESCRIPTION":  task.Description,
				"START-TIME":   fmt.Sprintf("%d", task.StartTimeInMillis),
				"RUNNING-TIME": fmt.Sprintf("%d", task.RunningTimeInNanos),
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
		sortCols := output.ParseSortColumns("NODE,ID")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
