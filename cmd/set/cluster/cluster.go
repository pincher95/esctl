package cluster

import "github.com/spf13/cobra"

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Update cluster-level settings",
}

func Cmd() *cobra.Command {
	clusterCmd.AddCommand(SettingsCmd())
	return clusterCmd
}
