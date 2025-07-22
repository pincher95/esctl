package cluster

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var clusterSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Print detailed information about an entity",
	Long:  "By default, it returns only settings that have been explicitly defined.",
	Example: utils.TrimAndIndent(`
	# Retrieve detailed information about the cluster settings.
	esctl describe cluster settings

	# Retrieve detailed information about the cluster settings in a none flat format.
	esctl describe cluster settings --no-flat-settings

	# Retrieve detailed information about the cluster settings including default settings.
	esctl describe cluster settings --include-defaults`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		handleDescribeClusterSettings(ctx)
		return nil
	},
}

func init() {
	clusterSettingsCmd.Flags().BoolVar(&flagFlatSettings, "no-flat-setting", false, "If set, print settings in a non-flat format (Default is false)")
	clusterSettingsCmd.Flags().BoolVar(&flagIncludeDefaults, "include-defaults", false, "If set, include default settings (Default is false)")
}

func handleDescribeClusterSettings(ctx context.Context) {
	settings, err := cluster.ClusterSettings(ctx, flagFlatSettings, flagIncludeDefaults)
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	output.PrintJson(settings)
}
