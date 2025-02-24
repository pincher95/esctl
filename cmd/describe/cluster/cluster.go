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
	- settings: Print detailed information about the cluster settings.`),
	Example: utils.TrimAndIndent(`
# Retrieve detailed information about the cluster settings.
esctl describe cluster settings

# Retrieve detailed information about the cluster settings in a flat format.
esctl describe cluster settings --flat-settings

# Retrieve detailed information about the cluster settings including default settings.
esctl describe cluster settings --include-defaults`),
}

func Cmd() *cobra.Command {
	return clusterCmd
}

func init() {
	clusterCmd.AddCommand(clusterSettingsCmd)
	clusterCmd.AddCommand(clusterHealthCmd)
}
