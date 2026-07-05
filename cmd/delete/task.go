package delete

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/tasks"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	deleteTaskID      string
	deleteTaskActions []string
)

var deleteTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Cancel a running task (or tasks by action)",
	Long: utils.Trim(`
Cancel one or more running tasks. Provide --task-id to cancel a specific task, or --actions to
cancel all tasks matching the given action patterns. Cancelling a task is a delete operation.`),
	Example: utils.TrimAndIndent(`
# Cancel a specific task
esctl delete task --task-id nodeId:12345

# Cancel all reindex tasks
esctl delete task --actions "*reindex"
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client := tasks.NewTasks()
		resp, err := client.CancelTasks(ctx, deleteTaskID, deleteTaskActions)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	deleteTaskCmd.Flags().StringArrayVarP(&deleteTaskActions, "actions", "a", nil, "Cancel tasks matching these action patterns")
	deleteTaskCmd.Flags().StringVar(&deleteTaskID, "task-id", "", "Cancel a specific task by id")
}
