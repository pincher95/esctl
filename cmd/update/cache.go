package update

import (
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagExpandWildcards string
	flagAllowNoIndices  bool
	flagFielddata       bool
	flagCacheIndex      string
	flagCacheQuery      bool
	flagCacheRequest    bool
	flagCacheFields     string
)

var updateCacheClearCmd = &cobra.Command{
	Use:   "cache",
	Short: "Clears the cache in the cluster or in a specific index",
	Long: utils.Trim(`
	Clear the cache of one or more indices. For data streams, the API clears the caches of the stream's backing indices.
	By default, the clear cache API clears all caches. To clear only specific caches, use the fielddata, query, or request parameters. To clear the cache only of specific fields, use the fields parameter.
	`),
	Example: utils.TrimAndIndent(`
	# Clear all caches for a specific index.
	esctl update cache --index my_index

	# Clear only the fielddata cache for an index.
	esctl update cache --index my_index --fielddata

	# Clear the query and request caches for all indices.
	esctl update cache --query --request

	# Clear the fielddata cache for specific fields only.
	esctl update cache --index my_index --fields user_id,session_id
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		idxClient := index.NewIndex()
		return handleCacheLogic(cmd, idxClient)
	},
}

func init() {
	updateCacheClearCmd.Flags().StringVarP(&flagCacheIndex, "index", "i", "", "Index (or comma-separated indices/patterns) to clear; empty clears all")
	updateCacheClearCmd.Flags().BoolVar(&flagFielddata, "fielddata", false, "Clear the fielddata cache")
	updateCacheClearCmd.Flags().BoolVar(&flagCacheQuery, "query", false, "Clear the query cache")
	updateCacheClearCmd.Flags().BoolVar(&flagCacheRequest, "request", false, "Clear the request cache")
	updateCacheClearCmd.Flags().StringVar(&flagCacheFields, "fields", "", "Clear the fielddata cache of specific fields only (comma-separated)")
	updateCacheClearCmd.Flags().BoolVar(&flagAllowNoIndices, "allow-no-indices", false, "If false, error when a wildcard/alias/_all targets only missing or closed indices")
	updateCacheClearCmd.Flags().StringVar(&flagExpandWildcards, "expand-wildcards", "", "Wildcard expansion: comma-separated subset of all,open,closed,hidden,none")
}

// handleCacheLogic forwards only the flags the user explicitly set, so a bare
// 'esctl update cache' clears all caches (the Elasticsearch default) while any
// specified cache/target narrows the request.
func handleCacheLogic(cmd *cobra.Command, client index.Index) error {
	ctx := cmd.Context()

	params := map[string]string{}
	if cmd.Flags().Changed("fielddata") {
		params["fielddata"] = strconv.FormatBool(flagFielddata)
	}
	if cmd.Flags().Changed("query") {
		params["query"] = strconv.FormatBool(flagCacheQuery)
	}
	if cmd.Flags().Changed("request") {
		params["request"] = strconv.FormatBool(flagCacheRequest)
	}
	if cmd.Flags().Changed("fields") {
		params["fields"] = flagCacheFields
	}
	if cmd.Flags().Changed("expand-wildcards") {
		params["expand_wildcards"] = flagExpandWildcards
	}
	if cmd.Flags().Changed("allow-no-indices") {
		params["allow_no_indices"] = strconv.FormatBool(flagAllowNoIndices)
	}

	cache, err := client.CacheClear(ctx, flagCacheIndex, params)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	return output.Render(cache)
}
