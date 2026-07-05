package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/slm"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var flagSLMPolicyName string

var getSlmPoliciesCmd = &cobra.Command{
	Use:   "slm-policies",
	Short: "List Snapshot Lifecycle Management policies or get a specific one",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# List all SLM policies
		esctl get slm-policies

		# Get a specific policy
		esctl get slm-policies --name daily-snapshots -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagSLMPolicyName != "" {
			return handleGetSLMPolicy(ctx, flagSLMPolicyName)
		}
		return handleListSLMPolicies(ctx)
	},
}

func init() {
	getSlmPoliciesCmd.Flags().StringVar(&flagSLMPolicyName, "name", "", "Policy name (omit to list all)")
	getCmd.AddCommand(getSlmPoliciesCmd)
}

var slmColumns = []output.ColumnDefaults{
	{Header: "NAME", Type: output.Text},
	{Header: "SCHEDULE", Type: output.Text},
	{Header: "REPOSITORY", Type: output.Text},
	{Header: "NEXT-EXECUTION", Type: output.Date},
	{Header: "VERSION", Type: output.Number},
}

func handleListSLMPolicies(ctx context.Context) error {
	policies, err := slm.List(ctx)
	if err != nil {
		return err
	}
	if len(policies) == 0 {
		fmt.Println("No SLM policies found")
		return nil
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(policies)
	}

	data := make([][]string, 0, len(policies))
	for name, policy := range policies {
		data = append(data, []string{
			name,
			policy.Policy.Schedule,
			policy.Policy.Repository,
			policy.NextExecution,
			fmt.Sprintf("%d", policy.Version),
		})
	}

	sortBy := flagSortBy
	if sortBy == "" {
		sortBy = "NAME"
	}
	return output.PrintTable(slmColumns, data, output.ParseSortColumns(sortBy))
}

func handleGetSLMPolicy(ctx context.Context, name string) error {
	policy, err := slm.Get(ctx, name)
	if err != nil {
		return err
	}
	return output.Render(policy)
}
