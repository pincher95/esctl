package cluster

import (
	"encoding/json"
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagBody == "" {
			return fmt.Errorf("--body JSON is required")
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(flagBody), &body); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
		resp, err := cluster.UpdateClusterSettings(ctx, body)
		if err != nil {
			return err
		}
		output.PrintJson(resp)
		return nil
	},
}

func Cmd() *cobra.Command {
	setClusterSettingsCmd.Flags().StringVar(&flagBody, "body", "", "Raw JSON settings body")
	_ = setClusterSettingsCmd.MarkFlagRequired("body")
	return setClusterSettingsCmd
}
