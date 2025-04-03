package update

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagExpandWildcards string
	flagAllowNoIndices  bool
	flagFielddata       bool
)

var updateCacheClearCmd = &cobra.Command{
	Use:   "cache",
	Short: "Clears the cache in the cluster or in a specific index",
	Long: utils.Trim(`
	Clear the cache of one or more indices. For data streams, the API clears the caches of the stream's backing indices.
	By default, the clear cache API clears all caches. To clear only specific caches, use the fielddata, query, or request parameters. To clear the cache only of specific fields, use the fields parameter.
	`),
	Example: utils.TrimAndIndent(`
	# Update the number of replicas for an index.
	esctl update cache --index my_index --fielddata true
	`),
	Run: func(cmd *cobra.Command, args []string) {
		index := index.NewIndex()
		handleCacheLogic(index)
	},
}

func init() {
	updateCacheClearCmd.Flags().BoolVar(&flagAllowNoIndices, "allow-no-indices", false, "If false, the request returns an error if any wildcard expression, index alias, or _all value targets only missing or closed indices. This behavior applies even if the request targets other open indices.")
	updateCacheClearCmd.Flags().StringVar(&flagExpandWildcards, "expand-wildcards", "all", "Type of index that wildcard patterns can match. If the request can target data streams, this argument determines whether wildcard expressions match hidden data streams. Supports comma-separated values, such as open,hidden. Valid values are: all, open, closed, hidden, none.")
	updateCacheClearCmd.Flags().BoolVar(&flagFielddata, "fielddata", true, "If true, clears the fields cache. Use the fields parameter to clear the cache of specific fields only..")
}

func handleCacheLogic(client index.Index) {
	cache, err := client.CacheClear(nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve cache:", err)
		os.Exit(1)
	}

	output.PrintJson(cache)
}
