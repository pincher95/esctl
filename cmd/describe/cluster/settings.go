package cluster

import (
	"fmt"

	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var clusterSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Print detailed information about an entity",
	Run: func(cmd *cobra.Command, args []string) {
		handleDescribeClusterSettings()
	},
}

func handleDescribeClusterSettings() {
	cluster, err := es.GetClusterSettings(flagFlatSettings, flagIncludeDefaults)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(cluster)
}

func init() {
	clusterSettingsCmd.Flags().BoolVar(&flagFlatSettings, "flat-settings", true, "If set, print settings in a flat format")
	clusterSettingsCmd.Flags().BoolVar(&flagIncludeDefaults, "include-defaults", false, "If set, include default settings")
}
