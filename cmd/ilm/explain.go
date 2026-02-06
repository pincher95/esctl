package ilm

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain <index>",
	Short: "Explain ILM status for indices",
	Args:  cobra.ExactArgs(1),
	Long: `Explain the ILM status for one or more indices.

This command shows the current ILM phase, action, and step for each index,
along with timing information and any failure details.`,
	Example: utils.TrimAndIndent(`
		# Explain ILM status for an index
		esctl ilm explain myindex

		# Explain for multiple indices using wildcard
		esctl ilm explain "logs-*"

		# Output as JSON
		esctl ilm explain myindex -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		index := args[0]

		result, err := ilm.Explain(ctx, index)
		if err != nil {
			return err
		}

		if len(result.Indices) == 0 {
			fmt.Println("No indices found matching pattern")
			return nil
		}

		// For JSON/YAML output
		if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
			return output.Render(result)
		}

		// For table output
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

		var data [][]string
		for indexName, info := range result.Indices {
			managed := "false"
			if info.Managed {
				managed = "true"
			}

			data = append(data, []string{
				indexName,
				managed,
				info.Policy,
				info.Phase,
				info.Action,
				info.Step,
				info.Age,
				info.FailedStep,
			})
		}

		return output.PrintTable(columnDefs, data, nil)
	},
}
