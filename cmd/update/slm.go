package update

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/slm"
	"github.com/spf13/cobra"
)

var updateSLMExecuteName string

var updateSlmExecuteCmd = &cobra.Command{
	Use:   "slm-execute",
	Short: "Trigger an SLM policy immediately (take a snapshot now)",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Run the policy now, outside its schedule
		esctl update slm-execute --name daily-snapshots
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := slm.Execute(cmd.Context(), updateSLMExecuteName)
		if err != nil {
			return err
		}
		fmt.Printf("SLM policy '%s' executed; snapshot '%s' started\n", updateSLMExecuteName, resp.SnapshotName)
		return nil
	},
}

func init() {
	updateSlmExecuteCmd.Flags().StringVar(&updateSLMExecuteName, "name", "", "Policy name (required)")
	_ = updateSlmExecuteCmd.MarkFlagRequired("name")
	updateCmd.AddCommand(updateSlmExecuteCmd)
}
