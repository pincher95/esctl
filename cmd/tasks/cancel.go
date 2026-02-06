package tasks

import (
	"github.com/pincher95/esctl/es/tasks"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagActions []string
	flagTaskID  string
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel running tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client := tasks.NewTasks()
		resp, err := client.CancelTasks(ctx, flagTaskID, flagActions)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	cancelCmd.Flags().StringArrayVarP(&flagActions, "actions", "a", nil, "Filter tasks by actions")
	cancelCmd.Flags().StringVar(&flagTaskID, "task-id", "", "Filter tasks by task ID")
}
