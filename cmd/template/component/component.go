package component

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "component",
	Short: "Manage component templates",
	Long: `Manage Elasticsearch component templates.

Component templates are building blocks for constructing index templates that specify index
settings, mappings, and aliases. They can be reused across multiple index templates to promote
consistency and reduce duplication.`,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(putCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(existsCmd)
}
