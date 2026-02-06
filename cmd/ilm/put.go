package ilm

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/internal/logger"
	"github.com/spf13/cobra"
)

var (
	flagFile string
)

var putCmd = &cobra.Command{
	Use:   "put <policy>",
	Short: "Create or update an ILM policy",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Create ILM policy from file
		esctl ilm put hot_delete_policy --file policy.json

		# Example policy file (policy.json):
		# {
		#   "policy": {
		#     "phases": {
		#       "hot": {
		#         "actions": {
		#           "rollover": {
		#             "max_age": "7d",
		#             "max_size": "50gb"
		#           }
		#         }
		#       },
		#       "delete": {
		#         "min_age": "30d",
		#         "actions": {
		#           "delete": {}
		#         }
		#       }
		#     }
		#   }
		# }
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		if flagFile == "" {
			return fmt.Errorf("--file is required")
		}

		logger.Debug("reading ILM policy from file", "file", flagFile)
		data, err := os.ReadFile(flagFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		policy, err := ilm.ParsePolicyFromJSON(data)
		if err != nil {
			return err
		}

		if err := ilm.Put(ctx, name, *policy); err != nil {
			return err
		}

		fmt.Printf("ILM policy '%s' created/updated successfully\n", name)
		return nil
	},
}

func init() {
	putCmd.Flags().StringVar(&flagFile, "file", "", "JSON file containing policy definition (required)")
	putCmd.MarkFlagRequired("file")
}
