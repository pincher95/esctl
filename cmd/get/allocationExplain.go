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
	Use:                   "explain [--include-disk-info] [--include-yes-decisions]",
	DisableFlagsInUseLine: true,
	Short:                 "Explain the shard allocations for the Elasticsearch cluster",
	Long: utils.Trim(`
	Get explanations for shard allocations in the cluster. For unassigned shards, it provides an explanation for why the shard is unassigned. For assigned shards, it provides an explanation for why the shard is remaining on its current node and has not moved or rebalanced to another node. This API can be very useful when attempting to diagnose why a shard is unassigned or why a shard continues to remain on its current node when you might expect otherwise.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all allocation explain.
	esctl get explain

	# Retrieve allocation explain with disk info.
	esctl get explain --include-disk-info

	# Retrieve allocation explain with yes decisions.
	esctl get explain --include-yes-decisions
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
