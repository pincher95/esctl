package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	repoType     string
	repoSettings string
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage snapshot repositories",
	Long:  "List, create, get, and delete snapshot repositories",
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snapshot repositories",
	Example: utils.TrimAndIndent(`
	# List all repositories
	esctl snapshot repo list
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleRepoList(ctx)
	},
}

var repoGetCmd = &cobra.Command{
	Use:   "get <repository>",
	Short: "Get details of a specific repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get repository details
	esctl snapshot repo get my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleRepoGet(ctx, args[0])
	},
}

var repoCreateCmd = &cobra.Command{
	Use:   "create <repository>",
	Short: "Create a new snapshot repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Create a filesystem repository
	esctl snapshot repo create my-repo --type=fs --settings="location:/backup/snapshots"

	# Create an S3 repository
	esctl snapshot repo create s3-repo --type=s3 --settings="bucket:my-bucket,base_path:snapshots"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleRepoCreate(ctx, args[0])
	},
}

var repoDeleteCmd = &cobra.Command{
	Use:   "delete <repository>",
	Short: "Delete a snapshot repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Delete a repository
	esctl snapshot repo delete my-repo
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleRepoDelete(ctx, args[0])
	},
}

func init() {
	repoCreateCmd.Flags().StringVar(&repoType, "type", "fs", "Repository type (fs, s3, azure, gcs, etc.)")
	repoCreateCmd.Flags().StringVar(&repoSettings, "settings", "", "Repository settings as key:value pairs separated by commas")
	repoCreateCmd.MarkFlagRequired("settings")

	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoGetCmd)
	repoCmd.AddCommand(repoCreateCmd)
	repoCmd.AddCommand(repoDeleteCmd)
}

func handleRepoList(ctx context.Context) error {
	repos, err := snapshots.ListRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	output.PrintJson(repos)
	return nil
}

func handleRepoGet(ctx context.Context, repository string) error {
	repo, err := snapshots.GetRepository(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	output.PrintJson(repo)
	return nil
}

func handleRepoCreate(ctx context.Context, repository string) error {
	settings, err := parseSettings(repoSettings)
	if err != nil {
		return fmt.Errorf("invalid settings format: %w", err)
	}

	err = snapshots.CreateRepository(ctx, repository, repoType, settings)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	fmt.Printf("Repository '%s' created successfully\n", repository)
	return nil
}

func handleRepoDelete(ctx context.Context, repository string) error {
	err := snapshots.DeleteRepository(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	fmt.Printf("Repository '%s' deleted successfully\n", repository)
	return nil
}

func parseSettings(settingsStr string) (map[string]interface{}, error) {
	settings := make(map[string]interface{})
	if settingsStr == "" {
		return settings, nil
	}

	pairs := strings.Split(settingsStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key:value pair: %s", pair)
		}
		settings[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	return settings, nil
}
