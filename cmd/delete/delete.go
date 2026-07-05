package delete

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Elasticsearch resources",
	Long:  "Delete Elasticsearch resources (indices, pipelines, snapshots, users, roles, etc.)",
}

func Cmd() *cobra.Command {
	deleteCmd.AddCommand(indexCmd)
	deleteCmd.AddCommand(byQueryCmd)
	deleteCmd.AddCommand(deleteAliasCmd)
	deleteCmd.AddCommand(deletePipelineCmd)
	deleteCmd.AddCommand(deleteSnapshotCmd)
	deleteCmd.AddCommand(deleteSnapshotRepoCmd)
	deleteCmd.AddCommand(deleteUserCmd)
	deleteCmd.AddCommand(deleteRoleCmd)
	deleteCmd.AddCommand(deleteReindexCmd)
	deleteCmd.AddCommand(dataStreamCmd)
	deleteCmd.AddCommand(deleteScriptCmd)
	deleteCmd.AddCommand(deleteSearchTemplateCmd)
	deleteCmd.AddCommand(deleteIlmPolicyCmd)
	deleteCmd.AddCommand(deleteSlmPolicyCmd)
	deleteCmd.AddCommand(deleteTaskCmd)
	deleteCmd.AddCommand(deleteTemplateCmd)
	deleteCmd.AddCommand(deleteComponentTemplateCmd)
	return deleteCmd
}
