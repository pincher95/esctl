package alias

import (
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage Elasticsearch index aliases",
	Long:  "Add, remove, and move aliases for Elasticsearch indices",
}

func Cmd() *cobra.Command {
	aliasCmd.AddCommand(addCmd)
	aliasCmd.AddCommand(removeCmd)
	aliasCmd.AddCommand(moveCmd)
	return aliasCmd
}
