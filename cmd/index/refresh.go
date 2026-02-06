package index

import (
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndex string
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh index (or all indices)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idx := index.NewIndex()
		resp, err := idx.Refresh(ctx, flagIndex)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	refreshCmd.Flags().StringVar(&flagIndex, "index", "", "Index name (empty means all indices)")
}
