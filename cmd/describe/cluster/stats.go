package cluster

import (
	"context"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		handleDescribeClusterStats(ctx)
		return nil
	},
}

func init() {
	clusterStatsCmd.Flags().BoolVar(&flagIncludeRemotes, "include-remotes", false, "Include remote cluster data into the response.")
	clusterStatsCmd.Flags().StringVar(&flagNodeID, "node-id", "", "Comma-separated list of node filters used to limit returned information. Defaults to all nodes in the cluster.")
}

func handleDescribeClusterStats(ctx context.Context) {
	stats, err := cluster.ClusterStats(ctx, flagNodeID, flagIncludeRemotes)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(stats)
}
