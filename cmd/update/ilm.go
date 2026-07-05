package update

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/spf13/cobra"
)

var updateILMRetryIndex string

var updateIlmRetryCmd = &cobra.Command{
	Use:   "ilm-retry",
	Short: "Retry the failed ILM step for an index",
	Long: utils.Trim(`
When an ILM action fails (rollover, shrink, force merge, ...), the index stays stuck on that step.
This retries the failed step. Use 'esctl get ilm-explain --index <name>' to inspect the status first.`),
	Args: cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Retry the failed ILM step for an index
		esctl update ilm-retry --index myindex

		# Retry for multiple indices via wildcard
		esctl update ilm-retry --index "logs-*"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		result, err := ilm.Explain(ctx, updateILMRetryIndex)
		if err != nil {
			return err
		}

		hasFailure := false
		for _, info := range result.Indices {
			if info.FailedStep != "" {
				hasFailure = true
				break
			}
		}
		if !hasFailure {
			fmt.Println("No failed steps found for the specified index(es)")
			return nil
		}

		if err := ilm.Retry(ctx, updateILMRetryIndex); err != nil {
			return err
		}

		fmt.Printf("ILM retry initiated for index '%s'\n", updateILMRetryIndex)
		fmt.Println("Use 'esctl get ilm-explain' to check the current status")
		return nil
	},
}

func init() {
	updateIlmRetryCmd.Flags().StringVar(&updateILMRetryIndex, "index", "", "Index name or pattern (required)")
	_ = updateIlmRetryCmd.MarkFlagRequired("index")
	updateCmd.AddCommand(updateIlmRetryCmd)
}
