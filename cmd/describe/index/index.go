package index

import (
	idx "github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagMappings        bool
	flagSettings        bool
	flagNoFlat          bool
	flagIncludeDefaults bool
)

var describeIndexCmd = &cobra.Command{
	Use:   "index [NAME]",
	Short: "Describe an index (mappings, settings)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		ctx := cmd.Context()
		// default: both if neither flag set
		m := flagMappings || (!flagMappings && !flagSettings)
		s := flagSettings || (!flagMappings && !flagSettings)
		resp, err := idx.GetIndexDetails(ctx, name, m, s, !flagNoFlat, flagIncludeDefaults)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func Cmd() *cobra.Command {
	describeIndexCmd.Flags().BoolVar(&flagMappings, "mappings", false, "Include mappings")
	describeIndexCmd.Flags().BoolVar(&flagSettings, "settings", true, "Include settings")
	describeIndexCmd.Flags().BoolVar(&flagNoFlat, "no-flat-setting", false, "Return nested (non-flat) settings in response")
	describeIndexCmd.Flags().BoolVar(&flagIncludeDefaults, "include-defaults", false, "Include default settings in response")
	return describeIndexCmd
}
