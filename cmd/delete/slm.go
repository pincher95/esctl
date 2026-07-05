package delete

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/slm"
	"github.com/spf13/cobra"
)

var (
	deleteSLMPolicyName string
	deleteSLMForce      bool
)

var deleteSlmPolicyCmd = &cobra.Command{
	Use:   "slm-policy",
	Short: "Delete a Snapshot Lifecycle Management policy",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Delete a policy (with confirmation)
		esctl delete slm-policy --name daily-snapshots

		# Delete without confirmation
		esctl delete slm-policy --name daily-snapshots --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if !deleteSLMForce {
			fmt.Printf("Delete SLM policy '%s'?\n", deleteSLMPolicyName)
			approved, err := utils.GetApproval()
			if err != nil {
				return err
			}
			if !approved {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := slm.Delete(ctx, deleteSLMPolicyName); err != nil {
			return err
		}

		fmt.Printf("SLM policy '%s' deleted successfully\n", deleteSLMPolicyName)
		return nil
	},
}

func init() {
	deleteSlmPolicyCmd.Flags().StringVar(&deleteSLMPolicyName, "name", "", "Policy name (required)")
	deleteSlmPolicyCmd.Flags().BoolVar(&deleteSLMForce, "force", false, "Skip confirmation")
	_ = deleteSlmPolicyCmd.MarkFlagRequired("name")
}
