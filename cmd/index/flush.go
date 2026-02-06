package index

import (
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush index (or all indices)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idx := index.NewIndex()
		resp, err := idx.Flush(ctx, flagIndex)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	flushCmd.Flags().StringVar(&flagIndex, "index", "", "Index name (empty means all indices)")
}
