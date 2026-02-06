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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ILM policies",
	Example: utils.TrimAndIndent(`
		# List all ILM policies
		esctl ilm list

		# Output as JSON
		esctl ilm list -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		policies, err := ilm.List(ctx)
		if err != nil {
			return err
		}

		if len(policies) == 0 {
			fmt.Println("No ILM policies found")
			return nil
		}

		// For JSON/YAML output
		if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
			return output.Render(policies)
		}

		// For table output
		columnDefs := []output.ColumnDefaults{
			{Header: "NAME", Type: output.Text},
			{Header: "VERSION", Type: output.Number},
			{Header: "PHASES", Type: output.Text},
			{Header: "MODIFIED", Type: output.Date},
		}

		var data [][]string
		for name, policy := range policies {
			phases := []string{}
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
	},
}

var flagSortBy string

func init() {
	listCmd.Flags().StringVarP(&flagSortBy, "sort-by", "s", "", "Column to sort by (NAME, VERSION, MODIFIED)")
}
