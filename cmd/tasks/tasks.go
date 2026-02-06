package tasks

import "github.com/spf13/cobra"

// Cmd returns the parent "tasks" command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Manage tasks",
	}
	cmd.AddCommand(cancelCmd)
	return cmd
}
