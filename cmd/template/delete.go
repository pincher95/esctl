package template

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
)

var (
	flagForce bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete an index template",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Delete template (with confirmation)
		esctl template delete logs-template

		# Delete without confirmation
		esctl template delete logs-template --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		// Check if template exists
		exists, err := template.Exists(ctx, name)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("template not found: %s", name)
		}

		// Confirm deletion
		if !flagForce {
			fmt.Printf("Delete template '%s'?\n", name)
			approved, err := utils.GetApproval()
			if err != nil {
				return err
			}
			if !approved {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := template.Delete(ctx, name); err != nil {
			return err
		}

		fmt.Printf("Template '%s' deleted successfully\n", name)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVar(&flagForce, "force", false, "Skip confirmation")
}
