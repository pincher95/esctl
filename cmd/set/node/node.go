package node

import "github.com/spf13/cobra"

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Update node-level settings",
}

func Cmd() *cobra.Command {
	nodeCmd.AddCommand(SettingsCmd())
	return nodeCmd
}
