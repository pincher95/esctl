package set

import (
	clustersettings "github.com/pincher95/esctl/cmd/set/cluster"
	indexsettings "github.com/pincher95/esctl/cmd/set/index"
	nodesettings "github.com/pincher95/esctl/cmd/set/node"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Create or update Elasticsearch resources",
}

func Cmd() *cobra.Command {
	setCmd.AddCommand(clustersettings.Cmd())
	setCmd.AddCommand(indexsettings.Cmd())
	setCmd.AddCommand(nodesettings.Cmd())
	setCmd.AddCommand(setAliasCmd)
	setCmd.AddCommand(setPipelineCmd)
	setCmd.AddCommand(setReindexCmd)
	setCmd.AddCommand(setSnapshotCmd)
	setCmd.AddCommand(setSnapshotRepoCmd)
	setCmd.AddCommand(setUserCmd)
	setCmd.AddCommand(setRoleCmd)
	setCmd.AddCommand(setScriptCmd)
	setCmd.AddCommand(setSearchTemplateCmd)
	setCmd.AddCommand(setIlmPolicyCmd)
	setCmd.AddCommand(setSlmPolicyCmd)
	setCmd.AddCommand(setTemplateCmd)
	setCmd.AddCommand(setComponentTemplateCmd)
	return setCmd
}
