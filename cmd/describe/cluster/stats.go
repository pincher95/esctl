package cluster

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var clusterStatsCmd = &cobra.Command{
	Use:     "stats",
	Short:   "Print detailed information about an entity",
	Long:    "By default, it returns only settings that have been explicitly defined.",
	Example: utils.TrimAndIndent(``),
	Run: func(cmd *cobra.Command, args []string) {
		handleDescribeClusterStats()
	},
}

func init() {
	clusterStatsCmd.Flags().BoolVar(&flagFlatSettings, "include-remotes", false, "Include remote cluster data into the response.")
	clusterStatsCmd.Flags().StringVar(&flagNodeID, "node-id", "", "Comma-separated list of node filters used to limit returned information. Defaults to all nodes in the cluster.")
}

func handleDescribeClusterStats() {
	stats, err := cluster.ClusterStats(nil, &flagNodeID, flagIncludeDefaults)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(stats)
}
