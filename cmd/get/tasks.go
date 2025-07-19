package get

import (
	"context"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/tasks"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Get tasks information",
	Long:  `This command retrieves and displays tasks information from Elasticsearch cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasksClient := tasks.NewTasks()
		cfg := config.ParseConfigFile()

		ctx := cmd.Context()

		if !flagRefresh {
			handleTaskLogic(ctx, tasksClient, *cfg)
			return
		}

		utils.WatchLoop(flagRefreshInterval, func() error {
			handleTaskLogic(ctx, tasksClient, *cfg)
			return nil
		})
	},
}

func init() {
	getTasksCmd.Flags().StringArrayVarP(&flagActions, "actions", "a", nil, "Filter tasks by actions")
	getTasksCmd.Flags().StringVar(&flagTasksID, "task-id", "", "Filter tasks by task ID")
}

var taskColumns = []output.ColumnDefaults{
	{Header: "NODE", Type: output.Text},
	{Header: "TASK-ID", Type: output.Text},
	{Header: "ID", Type: output.Number},
	{Header: "ACTION", Type: output.Text},
	{Header: "DESCRIPTION", Type: output.Text},
	{Header: "START-TIME", Type: output.Number},
	{Header: "RUNNING-TIME", Type: output.Number},
}

func handleTaskLogic(ctx context.Context, client tasks.Tasks, config config.Config) {
	tasksResponse, err := client.GetTasks(ctx, flagTasksID, flagActions)
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
				"NODE":         node.Name,
				"TASK-ID":      fmt.Sprintf("%s:%d", task.Node, task.ID),
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
		sortCols := output.ParseSortColumns("NODE,TASK-ID")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
