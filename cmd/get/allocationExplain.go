package get

import (
	"fmt"
	"time"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAllocationExplainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Get Elasticsearch allocation explain",
	Long: utils.Trim(`
	Get Elasticsearch allocation explain. You can filter the results using the node flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all allocation explain.
	esctl get explain

	# Retrieve allocation for a specific node.
	esctl get explain --node my_node
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		// config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			return handleAllocationExplainLogic()
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			err := handleAllocationExplainLogic()
			if err != nil {
				return err
			}
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getAllocationExplainCmd.Flags().BoolVar(&flagIncludeDiskInfo, "include-disk-info", false, "Information about disk usage and shard sizes")
	getAllocationExplainCmd.Flags().BoolVar(&flagIncludeYesDecisions, "include-yes-decisions", false, "YES decisions in explanation")
}

func handleAllocationExplainLogic() error {
	allocationsExplain, err := cluster.ClusterAllocationExplain(nil, flagIncludeDiskInfo, flagIncludeYesDecisions)
	if err != nil {
		return fmt.Errorf("Failed to retrieve allocation explain%v", err)
	}

	output.PrintJson(allocationsExplain)

	return nil
}
