package node

import (
	"fmt"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var describeNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Describe a node (cat nodes output)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := describeNodeFilter
		c := cat.NewCat()
		ctx := cmd.Context()
		nodes, err := c.CatNodes(ctx, "", filter, "", "")
		if err != nil {
			return err
		}
		if len(nodes) == 0 {
			return fmt.Errorf("node not found")
		}
		return output.Render(nodes)
	},
}

func Cmd() *cobra.Command {
	describeNodeCmd.Flags().StringVar(&describeNodeFilter, "name", "", "Filter by node name or substring of node name")
	return describeNodeCmd
}

var describeNodeFilter string
