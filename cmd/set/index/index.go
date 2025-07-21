package index

import "github.com/spf13/cobra"

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Update index-level settings",
}

func Cmd() *cobra.Command {
	indexCmd.AddCommand(SettingsCmd())
	return indexCmd
}
