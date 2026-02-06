package index

import "github.com/spf13/cobra"

// Cmd returns the parent "index" command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Index maintenance commands",
	}
	cmd.AddCommand(refreshCmd)
	cmd.AddCommand(flushCmd)
	cmd.AddCommand(forcemergeCmd)
	return cmd
}
