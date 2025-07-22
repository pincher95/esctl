package delete

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Elasticsearch resources",
	Long:  "Delete indices, documents by query, and aliases in Elasticsearch",
}

func Cmd() *cobra.Command {
	deleteCmd.AddCommand(indexCmd)
	deleteCmd.AddCommand(byQueryCmd)
	deleteCmd.AddCommand(aliasCmd)
	return deleteCmd
}
