package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var (
	flagILMPolicyName   string
	flagILMExplainIndex string
)

var getIlmPoliciesCmd = &cobra.Command{
	Use:   "ilm-policies",
	Short: "List ILM policies or get a specific policy",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# List all ILM policies
		esctl get ilm-policies

		# Get a specific policy
		esctl get ilm-policies --name hot_delete_policy

		# As JSON
		esctl get ilm-policies --name hot_delete_policy -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return runWithWatch(ctx, func() error {
			if flagILMPolicyName != "" {
				return handleGetILMPolicy(ctx, flagILMPolicyName)
			}
			return handleListILMPolicies(ctx)
		})
	},
}

var getIlmExplainCmd = &cobra.Command{
	Use:   "ilm-explain",
	Short: "Explain ILM status (phase/action/step) for indices",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Explain ILM status for an index
		esctl get ilm-explain --index myindex

		# Multiple indices via wildcard
		esctl get ilm-explain --index "logs-*"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleILMExplain(cmd.Context(), flagILMExplainIndex)
	},
}

func init() {
	getIlmPoliciesCmd.Flags().StringVar(&flagILMPolicyName, "name", "", "Policy name (omit to list all)")
	getIlmExplainCmd.Flags().StringVar(&flagILMExplainIndex, "index", "", "Index name or pattern")
	_ = getIlmExplainCmd.MarkFlagRequired("index")

	getCmd.AddCommand(getIlmPoliciesCmd)
	getCmd.AddCommand(getIlmExplainCmd)
}

func handleListILMPolicies(ctx context.Context) error {
	policies, err := ilm.List(ctx)
	if err != nil {
		return err
	}
	if len(policies) == 0 {
		fmt.Println("No ILM policies found")
		return nil
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(policies)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "NAME", Type: output.Text},
		{Header: "VERSION", Type: output.Number},
		{Header: "PHASES", Type: output.Text},
		{Header: "MODIFIED", Type: output.Date},
	}

	data := make([][]string, 0, len(policies))
	for name, policy := range policies {
		phases := make([]string, 0, len(policy.Policy.Phases))
		for phaseName := range policy.Policy.Phases {
			phases = append(phases, phaseName)
		}
		data = append(data, []string{
			name,
			fmt.Sprintf("%d", policy.Version),
			strings.Join(phases, ", "),
			policy.Modified,
		})
	}

	return output.PrintTable(columnDefs, data, output.ParseSortColumns(flagSortBy))
}

func handleGetILMPolicy(ctx context.Context, name string) error {
	policy, err := ilm.Get(ctx, name)
	if err != nil {
		return err
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(policy)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "setting", Type: output.Text},
		{Header: "value", Type: output.Text},
	}
	data := [][]string{
		{"name", policy.Name},
		{"version", fmt.Sprintf("%d", policy.Version)},
		{"modified", policy.Modified},
	}

	phases := make([]string, 0, len(policy.Policy.Phases))
	for phaseName := range policy.Policy.Phases {
		phases = append(phases, phaseName)
	}
	data = append(data, []string{"phases", strings.Join(phases, ", ")})

	for phaseName, phase := range policy.Policy.Phases {
		if phase.MinAge != "" {
			data = append(data, []string{phaseName + ".min_age", phase.MinAge})
		}
		actions := make([]string, 0, len(phase.Actions))
		for actionName := range phase.Actions {
			actions = append(actions, actionName)
		}
		data = append(data, []string{phaseName + ".actions", strings.Join(actions, ", ")})
	}

	return output.PrintTable(columnDefs, data, nil)
}

func handleILMExplain(ctx context.Context, index string) error {
	result, err := ilm.Explain(ctx, index)
	if err != nil {
		return err
	}
	if len(result.Indices) == 0 {
		fmt.Println("No indices found matching pattern")
		return nil
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(result)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "INDEX", Type: output.Text},
		{Header: "MANAGED", Type: output.Boolean},
		{Header: "POLICY", Type: output.Text},
		{Header: "PHASE", Type: output.Text},
		{Header: "ACTION", Type: output.Text},
		{Header: "STEP", Type: output.Text},
		{Header: "AGE", Type: output.Text},
		{Header: "FAILED-STEP", Type: output.Text},
	}

	data := make([][]string, 0, len(result.Indices))
	for indexName, info := range result.Indices {
		managed := "false"
		if info.Managed {
			managed = "true"
		}
		data = append(data, []string{
			indexName, managed, info.Policy, info.Phase,
			info.Action, info.Step, info.Age, info.FailedStep,
		})
	}

	return output.PrintTable(columnDefs, data, output.ParseSortColumns(flagSortBy))
}
