package delete

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/spf13/cobra"
)

var (
	deleteILMPolicyName string
	deleteILMForce      bool
)

var deleteIlmPolicyCmd = &cobra.Command{
	Use:   "ilm-policy",
	Short: "Delete an ILM policy",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Delete a policy (with confirmation)
		esctl delete ilm-policy --name old_policy

		# Delete without confirmation
		esctl delete ilm-policy --name old_policy --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		exists, err := ilm.Exists(ctx, deleteILMPolicyName)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("ILM policy not found: %s", deleteILMPolicyName)
		}

		if !deleteILMForce {
			fmt.Printf("Delete ILM policy '%s'?\n", deleteILMPolicyName)
			approved, err := utils.GetApproval()
			if err != nil {
				return err
			}
			if !approved {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := ilm.Delete(ctx, deleteILMPolicyName); err != nil {
			return err
		}

		fmt.Printf("ILM policy '%s' deleted successfully\n", deleteILMPolicyName)
		return nil
	},
}

func init() {
	deleteIlmPolicyCmd.Flags().StringVar(&deleteILMPolicyName, "name", "", "Policy name (required)")
	deleteIlmPolicyCmd.Flags().BoolVar(&deleteILMForce, "force", false, "Skip confirmation")
	_ = deleteIlmPolicyCmd.MarkFlagRequired("name")
}
