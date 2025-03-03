package cluster

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Print detailed information about an entity",
	Long: utils.Trim(`
The 'describe cluster' command allows you to retrieve detailed information about the Elasticsearch cluster.

Available Subcommands:
	- settings: Print detailed information about the cluster settings.
	- health: Print detailed information about the cluster health.`),
}

func Cmd() *cobra.Command {
	return clusterCmd
}

func init() {
	clusterCmd.AddCommand(clusterSettingsCmd)
	clusterCmd.AddCommand(clusterHealthCmd)
	clusterCmd.AddCommand(clusterStatsCmd)
}
