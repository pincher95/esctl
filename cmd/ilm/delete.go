package ilm

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/spf13/cobra"
)

var (
	flagForce bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <policy>",
	Short: "Delete an ILM policy",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Delete policy (with confirmation)
		esctl ilm delete old_policy

		# Delete without confirmation
		esctl ilm delete old_policy --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		// Check if policy exists
		exists, err := ilm.Exists(ctx, name)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("ILM policy not found: %s", name)
		}

		// Confirm deletion
		if !flagForce {
			fmt.Printf("Delete ILM policy '%s'?\n", name)
			approved, err := utils.GetApproval()
			if err != nil {
				return err
			}
			if !approved {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := ilm.Delete(ctx, name); err != nil {
			return err
		}

		fmt.Printf("ILM policy '%s' deleted successfully\n", name)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&flagForce, "force", false, "Skip confirmation")
}
