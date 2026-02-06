package ilm

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "ilm",
	Short: "Manage Index Lifecycle Management policies",
	Long: `Manage ILM policies in Elasticsearch/OpenSearch.

Index Lifecycle Management (ILM) policies define how indices are managed over time,
including actions like rollover, shrink, force merge, and delete. ILM automates
the management of time-series indices based on age, size, and other criteria.`,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(putCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(existsCmd)
	Cmd.AddCommand(explainCmd)
	Cmd.AddCommand(retryCmd)
}
