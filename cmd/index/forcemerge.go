package index

import (
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagMaxNumSegments     int
	flagOnlyExpungeDeletes bool
	flagFlush              bool
)

var forcemergeCmd = &cobra.Command{
	Use:   "forcemerge",
	Short: "Force-merge index segments",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idx := index.NewIndex()
		resp, err := idx.Forcemerge(ctx, flagIndex, flagMaxNumSegments, flagOnlyExpungeDeletes, flagFlush)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	forcemergeCmd.Flags().StringVar(&flagIndex, "index", "", "Index name (empty means all indices)")
	forcemergeCmd.Flags().IntVar(&flagMaxNumSegments, "max-num-segments", 0, "Maximum number of segments to merge to")
	forcemergeCmd.Flags().BoolVar(&flagOnlyExpungeDeletes, "only-expunge-deletes", false, "Only expunge deleted documents")
	forcemergeCmd.Flags().BoolVar(&flagFlush, "flush", false, "Refresh the index after force-merge")
}
