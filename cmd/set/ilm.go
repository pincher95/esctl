package set

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/spf13/cobra"
)

var (
	setILMPolicyName string
	setILMPolicyFile string
)

var setIlmPolicyCmd = &cobra.Command{
	Use:   "ilm-policy",
	Short: "Create or update an ILM policy",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Create or update a policy from a file
		esctl set ilm-policy --name hot_delete_policy --file policy.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		data, err := os.ReadFile(setILMPolicyFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		policy, err := ilm.ParsePolicyFromJSON(data)
		if err != nil {
			return err
		}

		if err := ilm.Put(ctx, setILMPolicyName, *policy); err != nil {
			return err
		}

		fmt.Printf("ILM policy '%s' created/updated successfully\n", setILMPolicyName)
		return nil
	},
}

func init() {
	setIlmPolicyCmd.Flags().StringVar(&setILMPolicyName, "name", "", "Policy name (required)")
	setIlmPolicyCmd.Flags().StringVar(&setILMPolicyFile, "file", "", "JSON file containing the policy definition (required)")
	_ = setIlmPolicyCmd.MarkFlagRequired("name")
	_ = setIlmPolicyCmd.MarkFlagRequired("file")
}
