package component

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
)

var flagExistsName string

var existsCmd = &cobra.Command{
	Use:   "exists",
	Short: "Check if a component template exists",
	Long: utils.Trim(`
		Checks whether a component template exists in the Elasticsearch cluster.
	`),
	Example: utils.TrimAndIndent(`
		# Check if a component template exists
		esctl template component exists --name my-component
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleExistsComponent(ctx, flagExistsName)
	},
}

func init() {
	existsCmd.Flags().StringVar(&flagExistsName, "name", "", "Component template name")
	_ = existsCmd.MarkFlagRequired("name")
}

func handleExistsComponent(ctx context.Context, name string) error {
	exists, err := template.ExistsComponent(ctx, name)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Component template '%s' exists\n", name)
	} else {
		fmt.Printf("Component template '%s' does not exist\n", name)
	}

	return nil
}
