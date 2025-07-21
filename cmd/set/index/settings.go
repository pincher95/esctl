package index

import (
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndexName string
	flagBody      string
	flagFlat      bool
)

var setIndexSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Set index settings (transient) for a given index",
	Long: `Update LIVE index settings such as refresh_interval or number_of_replicas.

Example:
  # set replicas on my-index
  esctl set index settings --index my-index --body '{"number_of_replicas":2}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagIndexName == "" {
			return fmt.Errorf("--index is required")
		}
		if flagBody == "" {
			return fmt.Errorf("--body JSON is required")
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(flagBody), &body); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
		idx := index.NewIndex()
		resp, err := idx.UpdateIndexSettings(ctx, flagIndexName, body, flagFlat)
		if err != nil {
			return err
		}
		output.PrintJson(resp)
		return nil
	},
}

func SettingsCmd() *cobra.Command {
	setIndexSettingsCmd.Flags().StringVar(&flagIndexName, "index", "", "Target index name")
	_ = setIndexSettingsCmd.MarkFlagRequired("index")
	setIndexSettingsCmd.Flags().StringVar(&flagBody, "body", "", "Raw JSON settings body")
	_ = setIndexSettingsCmd.MarkFlagRequired("body")
	setIndexSettingsCmd.Flags().BoolVar(&flagFlat, "no-flat", false, "Return noneflat settings in response")
	return setIndexSettingsCmd
}
