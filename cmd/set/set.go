package set

import (
	clustersettings "github.com/pincher95/esctl/cmd/set/cluster"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Lightweight updates to cluster or index settings",
}

func Cmd() *cobra.Command {
	setCmd.AddCommand(clustersettings.Cmd())
	return setCmd
}
