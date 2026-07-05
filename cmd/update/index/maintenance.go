package index

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

// These commands read the parent "update index" persistent --index/-i flag via
// cmd.Flags().GetString("index"); an empty value targets all indices.

var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh an index (or all indices) to make recent changes searchable",
	Example: utils.TrimAndIndent(`
		esctl update index refresh --index my-index
		esctl update index refresh            # all indices
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		indexName, _ := cmd.Flags().GetString("index")
		resp, err := index.NewIndex().Refresh(cmd.Context(), indexName)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

var FlushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush an index (or all indices) to disk",
	Example: utils.TrimAndIndent(`
		esctl update index flush --index my-index
		esctl update index flush              # all indices
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		indexName, _ := cmd.Flags().GetString("index")
		resp, err := index.NewIndex().Flush(cmd.Context(), indexName)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

var (
	flagFMMaxNumSegments     int
	flagFMOnlyExpungeDeletes bool
	flagFMFlush              bool
)

var ForcemergeCmd = &cobra.Command{
	Use:   "forcemerge",
	Short: "Force-merge index segments",
	Long: utils.Trim(`
Merge the segments of an index to reduce their number. Force-merge is expensive; on read-only
indices merging to a small --max-num-segments (often 1) can significantly reduce heap and disk use.`),
	Example: utils.TrimAndIndent(`
		esctl update index forcemerge --index my-index --max-num-segments 1
		esctl update index forcemerge --index my-index --only-expunge-deletes
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		indexName, _ := cmd.Flags().GetString("index")
		resp, err := index.NewIndex().Forcemerge(cmd.Context(), indexName, flagFMMaxNumSegments, flagFMOnlyExpungeDeletes, flagFMFlush)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	ForcemergeCmd.Flags().IntVar(&flagFMMaxNumSegments, "max-num-segments", 0, "Maximum number of segments to merge to")
	ForcemergeCmd.Flags().BoolVar(&flagFMOnlyExpungeDeletes, "only-expunge-deletes", false, "Only expunge deleted documents")
	ForcemergeCmd.Flags().BoolVar(&flagFMFlush, "flush", false, "Flush the index after force-merge")
}
