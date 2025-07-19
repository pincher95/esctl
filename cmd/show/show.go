package show

import (
	"github.com/pincher95/esctl/cmd/get"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show Elasticsearch resources. Use -o json|yaml for structured output",
}

func Cmd() *cobra.Command {
	// Re-use all existing get sub-commands under show.
	showCmd.AddCommand(get.Cmd())
	return showCmd
}
