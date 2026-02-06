package ilm

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/spf13/cobra"
)

var existsCmd = &cobra.Command{
	Use:   "exists",
	Short: "Check if an ILM policy exists",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Check if policy exists
		esctl ilm exists --name hot_delete_policy
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := existsPolicyName

		exists, err := ilm.Exists(ctx, name)
		if err != nil {
			return err
		}

		if exists {
			fmt.Printf("ILM policy '%s' exists\n", name)
		} else {
			fmt.Printf("ILM policy '%s' does not exist\n", name)
		}

		return nil
	},
}

var existsPolicyName string

func init() {
	existsCmd.Flags().StringVar(&existsPolicyName, "name", "", "Policy name")
	existsCmd.MarkFlagRequired("name")
}
