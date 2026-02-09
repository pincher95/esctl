package cluster

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagBody string
)

var setClusterSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Set cluster settings (transient/persistent)",
	Long: `Update cluster‐wide settings without closing indices.

Either or both of the top-level objects "transient" or "persistent" must be supplied.

Examples:

  # Throttle recovery to 100 MB/s persistently
  esctl set settings --body '{"persistent":{"indices.recovery.max_bytes_per_sec":"100mb"}}'

  # Enable shard allocation temporarily (transient)
  esctl set settings --body '{"transient":{"cluster.routing.allocation.enable":"all"}}'`,
	Example: `  # disable rebalancing persistently
  esctl set settings --body '{"persistent":{"cluster.routing.rebalance.enable":"none"}}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagBody == "" {
			return fmt.Errorf("--body JSON is required")
		}
		body, err := utils.ParseJSONMap(flagBody, "invalid JSON")
		if err != nil {
			return err
		}
		resp, err := cluster.UpdateClusterSettings(ctx, body)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func SettingsCmd() *cobra.Command {
	setClusterSettingsCmd.Flags().StringVar(&flagBody, "body", "", "Raw JSON settings body")
	_ = setClusterSettingsCmd.MarkFlagRequired("body")
	return setClusterSettingsCmd
}
