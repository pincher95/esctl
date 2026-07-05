package set

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/slm"
	"github.com/spf13/cobra"
)

var (
	setSLMPolicyName string
	setSLMPolicyFile string
)

var setSlmPolicyCmd = &cobra.Command{
	Use:   "slm-policy",
	Short: "Create or update a Snapshot Lifecycle Management policy",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Create or update a policy from a file
		esctl set slm-policy --name daily-snapshots --file policy.json

		# Example policy.json:
		# {
		#   "schedule": "0 30 1 * * ?",
		#   "name": "<daily-snap-{now/d}>",
		#   "repository": "my_repository",
		#   "config": { "indices": ["*"] },
		#   "retention": { "expire_after": "30d", "min_count": 5, "max_count": 50 }
		# }
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		data, err := os.ReadFile(setSLMPolicyFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if err := slm.Put(ctx, setSLMPolicyName, data); err != nil {
			return err
		}

		fmt.Printf("SLM policy '%s' created/updated successfully\n", setSLMPolicyName)
		return nil
	},
}

func init() {
	setSlmPolicyCmd.Flags().StringVar(&setSLMPolicyName, "name", "", "Policy name (required)")
	setSlmPolicyCmd.Flags().StringVar(&setSLMPolicyFile, "file", "", "JSON file containing the policy definition (required)")
	_ = setSlmPolicyCmd.MarkFlagRequired("name")
	_ = setSlmPolicyCmd.MarkFlagRequired("file")
}
