package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
)

func HandleRepoList(ctx context.Context, nameFilter string) error {
	repos, err := snapshots.ListRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}
	if nameFilter != "" {
		filtered := make(snapshots.RepositoryResponse)
		for name, repo := range repos {
			if strings.Contains(name, nameFilter) {
				filtered[name] = repo
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("no repositories matched: %s", nameFilter)
		}
		return output.Render(filtered)
	}
	return output.Render(repos)
}

func HandleRepoGet(ctx context.Context, repository string) error {
	repo, err := snapshots.GetRepository(ctx, repository)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	return output.Render(repo)
}

func HandleRepoCreate(ctx context.Context, repository, repoTypeValue, settingsValue string) error {
	settings, err := ParseSettings(settingsValue)
	if err != nil {
		return fmt.Errorf("invalid settings format: %w", err)
	}

	if err := snapshots.CreateRepository(ctx, repository, repoTypeValue, settings); err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	fmt.Printf("Repository '%s' created successfully\n", repository)
	return nil
}

func HandleRepoDelete(ctx context.Context, repository string) error {
	if err := snapshots.DeleteRepository(ctx, repository); err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	fmt.Printf("Repository '%s' deleted successfully\n", repository)
	return nil
}

func ParseSettings(settingsStr string) (map[string]interface{}, error) {
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
