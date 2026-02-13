package delete

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/searchtemplate"
	"github.com/spf13/cobra"
)

var (
	flagDeleteSearchTemplateID     string
	flagDeleteSearchTemplateDryRun bool
)

var deleteSearchTemplateCmd = &cobra.Command{
	Use:   "search-template",
	Short: "Delete a stored search template",
	Long: utils.Trim(`
		Deletes a stored search template from Elasticsearch. This operation cannot be undone.
	`),
	Example: utils.TrimAndIndent(`
		# Delete a search template
		esctl delete search-template --id my-template

		# Dry run - show what would be deleted
		esctl delete search-template --id my-template --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleDeleteSearchTemplateLogic(ctx, flagDeleteSearchTemplateID, flagDeleteSearchTemplateDryRun)
	},
}

func init() {
	deleteSearchTemplateCmd.Flags().StringVar(&flagDeleteSearchTemplateID, "id", "", "Search template ID")
	deleteSearchTemplateCmd.Flags().BoolVar(&flagDeleteSearchTemplateDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	_ = deleteSearchTemplateCmd.MarkFlagRequired("id")
}

func handleDeleteSearchTemplateLogic(ctx context.Context, id string, dryRun bool) error {
	if dryRun {
		fmt.Printf("DRY RUN: Would delete search template: %s\n", id)
		return nil
	}

	if err := searchtemplate.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete search template: %w", err)
	}

	fmt.Printf("Successfully deleted search template: %s\n", id)
	return nil
}
