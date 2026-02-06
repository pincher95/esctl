package node

import (
	"fmt"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var describeNodeCmd = &cobra.Command{
	Use:   "node [NAME]",
	Short: "Describe a node (cat nodes output)",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var filter string
		if len(args) == 1 {
			filter = args[0]
		}
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

func Cmd() *cobra.Command { return describeNodeCmd }
