package template

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "template",
	Short: "Manage index templates",
	Long: `Manage Elasticsearch index templates.

Index templates let you define settings, mappings, and aliases that are automatically
applied when new indices are created that match the template's patterns.`,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(putCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(existsCmd)
}
