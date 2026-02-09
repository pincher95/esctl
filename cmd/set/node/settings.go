package node

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/node"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagNodeID string
	flagBody   string
	flagFlat   bool
)

var setNodeSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Set live settings for one or all nodes",
	Long: `Update node-level settings via the undocumented _nodes/<id>/settings endpoint.

Example:
  # Increase processors for node-1
  esctl set node settings --node-id node-1 --body '{"processors":4}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagBody == "" {
			return fmt.Errorf("--body JSON is required")
		}
		body, err := utils.ParseJSONMap(flagBody, "invalid JSON")
		if err != nil {
			return err
		}
		resp, err := node.UpdateNodeSettings(ctx, flagNodeID, body, flagFlat)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func SettingsCmd() *cobra.Command {
	setNodeSettingsCmd.Flags().StringVar(&flagNodeID, "node-id", "", "Target node ID (empty = all nodes)")
	setNodeSettingsCmd.Flags().StringVar(&flagBody, "body", "", "Raw JSON settings body")
	_ = setNodeSettingsCmd.MarkFlagRequired("body")
	setNodeSettingsCmd.Flags().BoolVar(&flagFlat, "no-flat", false, "Return none flat settings in response")
	return setNodeSettingsCmd
}
