package delete

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
)

var (
	deleteTemplateName          string
	deleteTemplateForce         bool
	deleteComponentTemplateName string
	deleteComponentForce        bool
)

var deleteTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Delete an index template",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		esctl delete template --name logs-template
		esctl delete template --name logs-template --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		exists, err := template.Exists(ctx, deleteTemplateName)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("template not found: %s", deleteTemplateName)
		}

		if !deleteTemplateForce {
			fmt.Printf("Delete template '%s'?\n", deleteTemplateName)
			approved, err := utils.GetApproval()
			if err != nil {
				return err
			}
			if !approved {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := template.Delete(ctx, deleteTemplateName); err != nil {
			return err
		}
		fmt.Printf("Template '%s' deleted successfully\n", deleteTemplateName)
		return nil
	},
}

var deleteComponentTemplateCmd = &cobra.Command{
	Use:   "component-template",
	Short: "Delete a component template",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		esctl delete component-template --name settings-common
		esctl delete component-template --name settings-common --force
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		exists, err := template.ExistsComponent(ctx, deleteComponentTemplateName)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("component template not found: %s", deleteComponentTemplateName)
		}

		if !deleteComponentForce {
			fmt.Printf("Delete component template '%s'?\n", deleteComponentTemplateName)
			approved, err := utils.GetApproval()
			if err != nil {
				return err
			}
			if !approved {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := template.DeleteComponent(ctx, deleteComponentTemplateName); err != nil {
			return err
		}
		fmt.Printf("Component template '%s' deleted successfully\n", deleteComponentTemplateName)
		return nil
	},
}

func init() {
	deleteTemplateCmd.Flags().StringVar(&deleteTemplateName, "name", "", "Template name (required)")
	deleteTemplateCmd.Flags().BoolVar(&deleteTemplateForce, "force", false, "Skip confirmation")
	_ = deleteTemplateCmd.MarkFlagRequired("name")

	deleteComponentTemplateCmd.Flags().StringVar(&deleteComponentTemplateName, "name", "", "Component template name (required)")
	deleteComponentTemplateCmd.Flags().BoolVar(&deleteComponentForce, "force", false, "Skip confirmation")
	_ = deleteComponentTemplateCmd.MarkFlagRequired("name")
}
