package delete

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/script"
	"github.com/spf13/cobra"
)

var (
	flagDeleteScriptID     string
	flagDeleteScriptDryRun bool
)

var deleteScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Delete a stored script",
	Long: utils.Trim(`
		Deletes a stored script from Elasticsearch. This operation cannot be undone.
	`),
	Example: utils.TrimAndIndent(`
		# Delete a script
		esctl delete script --id my-script

		# Dry run - show what would be deleted
		esctl delete script --id my-script --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleDeleteScriptLogic(ctx, flagDeleteScriptID, flagDeleteScriptDryRun)
	},
}

func init() {
	deleteScriptCmd.Flags().StringVar(&flagDeleteScriptID, "id", "", "Script ID")
	deleteScriptCmd.Flags().BoolVar(&flagDeleteScriptDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	_ = deleteScriptCmd.MarkFlagRequired("id")
}

func handleDeleteScriptLogic(ctx context.Context, id string, dryRun bool) error {
	if dryRun {
		fmt.Printf("DRY RUN: Would delete script: %s\n", id)
		return nil
	}

	if err := script.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete script: %w", err)
	}

	fmt.Printf("Successfully deleted script: %s\n", id)
	return nil
}
