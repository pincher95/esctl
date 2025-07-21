package set

import (
	clustersettings "github.com/pincher95/esctl/cmd/set/cluster"
	indexsettings "github.com/pincher95/esctl/cmd/set/index"
	nodesettings "github.com/pincher95/esctl/cmd/set/node"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Lightweight updates to cluster or index settings",
}

func Cmd() *cobra.Command {
	setCmd.AddCommand(clustersettings.Cmd())
	setCmd.AddCommand(indexsettings.Cmd())
	setCmd.AddCommand(nodesettings.Cmd())
	return setCmd
}
