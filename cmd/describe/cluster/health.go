package cluster

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var clusterHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get detailed information about the cluster health",
	Long:  "Can be one of cluster, indices or shards. Controls the details level of the health information returned.",
	Example: utils.TrimAndIndent(`
	# Retrieve detailed information about the cluster health.
	esctl describe cluster health

	# Retrieve detailed information about the cluster health in a flat format.
	esctl describe cluster health --level=indices

	# Retrieve detailed information about the cluster health including default settings.
	esctl describe cluster health --level=shards

	# Retrieve detailed information about the cluster health including default settings.
	esctl describe cluster health --index=my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		handleDescribeClusterHealth()
	},
}

func init() {
	clusterHealthCmd.Flags().StringVar(&flagLevel, "level", "", "If set, print settings in a flat format")
	clusterHealthCmd.Flags().StringVar(&flagIndex, "index", "", "Comma-separated list of data streams, indices, and index aliases used to limit the request.")
	clusterHealthCmd.Flags().StringVar(&flagExpandWildcards, "expand-wildcards", "", "Expands wildcard expressions to concrete indexes. Combine multiple values with commas. Supported values are all, open, closed, hidden, and none. (Default is open)")
}

func handleDescribeClusterHealth() {
	health, err := cluster.ClusterHealth(nil, &flagLevel, &flagExpandWildcards, &flagIndex)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(health)
}
