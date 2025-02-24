package cluster

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es"
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

func handleDescribeClusterHealth() {
	cluster, err := es.GetClusterHealth(flagLevel, flagIndex)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(cluster)
}

func init() {
	clusterHealthCmd.Flags().StringVar(&flagLevel, "level", "", "If set, print settings in a flat format")
	clusterHealthCmd.Flags().StringVar(&flagIndex, "index", "", "Comma-separated list of data streams, indices, and index aliases used to limit the request.")
}
