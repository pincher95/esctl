package ilm

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/spf13/cobra"
)

var retryCmd = &cobra.Command{
	Use:   "retry <index>",
	Short: "Retry failed ILM step for an index",
	Args:  cobra.ExactArgs(1),
	Long: `Retry the failed ILM step for an index.

When an ILM action fails (e.g., rollover, shrink, force merge), the index
remains stuck in that step. This command retries the failed step.`,
	Example: utils.TrimAndIndent(`
		# Retry failed ILM step
		esctl ilm retry myindex

		# Retry for multiple indices
		esctl ilm retry "logs-*"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		index := args[0]

		// First check if there's a failed step
		result, err := ilm.Explain(ctx, index)
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

		if err := ilm.Retry(ctx, index); err != nil {
			return err
		}

		fmt.Printf("ILM retry initiated for index '%s'\n", index)
		fmt.Println("Use 'esctl ilm explain' to check the current status")
		return nil
	},
}
