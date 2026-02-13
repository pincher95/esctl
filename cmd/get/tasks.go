package get

import (
	"context"
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		tasksClient := tasks.NewTasks()
		cfg, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		if !flagRefresh {
			return handleTaskLogic(ctx, tasksClient, *cfg)
		}

		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleTaskLogic(ctx, tasksClient, *cfg)
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

func handleTaskLogic(ctx context.Context, client tasks.Tasks, config config.Config) error {
	tasksResponse, err := client.GetTasks(ctx, flagTasksID, flagActions)
	if err != nil {
		return fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	columnDefs, err := getColumnDefs(config, "task", taskColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	// Calculate total capacity for pre-allocation
	totalTasks := 0
	for _, node := range tasksResponse.Nodes {
		totalTasks += len(node.Tasks)
	}
	data := make([][]string, 0, totalTasks)

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
		return output.PrintTable(columnDefs, data, sortCols)
	}
	sortCols := output.ParseSortColumns("NODE,TASK-ID")
	return output.PrintTable(columnDefs, data, sortCols)
}
