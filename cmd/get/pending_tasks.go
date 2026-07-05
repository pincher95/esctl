package get

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var getPendingTasksCmd = &cobra.Command{
	Use:   "pending-tasks",
	Short: "List cluster-state changes waiting to be applied by the master",
	Long: utils.Trim(`
Show tasks in the cluster-state update queue that have not yet been executed. A queue that stays
non-empty (especially with high time-in-queue) indicates an overloaded or stuck master node.`),
	Example: utils.TrimAndIndent(`
# Show the pending cluster tasks.
esctl get pending-tasks

# Watch the queue drain in real time.
esctl get pending-tasks --watch --interval 2s
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if !flagRefresh {
			return handlePendingTasksLogic(ctx)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handlePendingTasksLogic(ctx)
		})
	},
}

func init() {
	getCmd.AddCommand(getPendingTasksCmd)
}

var pendingTasksColumns = []output.ColumnDefaults{
	{Header: "INSERT-ORDER", Type: output.Number},
	{Header: "PRIORITY", Type: output.Text},
	{Header: "SOURCE", Type: output.Text},
	{Header: "TIME-IN-QUEUE", Type: output.Text},
	{Header: "EXECUTING", Type: output.Boolean},
}

func handlePendingTasksLogic(ctx context.Context) error {
	resp, err := cluster.ClusterPendingTasks(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve pending tasks: %w", err)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(resp)
	}

	data := make([][]string, 0, len(resp.Tasks))
	for _, task := range resp.Tasks {
		rowData := map[string]string{
			"INSERT-ORDER":  strconv.Itoa(task.InsertOrder),
			"PRIORITY":      task.Priority,
			"SOURCE":        task.Source,
			"TIME-IN-QUEUE": task.TimeInQueue,
			"EXECUTING":     strconv.FormatBool(task.Executing),
		}
		row := make([]string, len(pendingTasksColumns))
		for i, colDef := range pendingTasksColumns {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	sortBy := flagSortBy
	if sortBy == "" {
		sortBy = "INSERT-ORDER"
	}
	return output.PrintTable(pendingTasksColumns, data, output.ParseSortColumns(sortBy))
}
