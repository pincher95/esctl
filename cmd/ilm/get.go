package ilm

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific ILM policy",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Get ILM policy details
		esctl ilm get --name hot_delete_policy

		# Get as JSON
		esctl ilm get --name hot_delete_policy -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := getPolicyName

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
	},
}

var getPolicyName string

func init() {
	getCmd.Flags().StringVar(&getPolicyName, "name", "", "Policy name")
	getCmd.MarkFlagRequired("name")
}
