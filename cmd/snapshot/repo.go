package snapshot

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
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
		repos = filtered
	}
	columnDefs := []output.ColumnDefaults{
		{Header: "name", Type: output.Text},
		{Header: "type", Type: output.Text},
	}

	settingKeys := make(map[string]struct{})
	flatSettings := make(map[string]map[string]any, len(repos))
	for name, repo := range repos {
		flat := utils.FlattenSettingsMap(repo.Settings)
		flatSettings[name] = flat
		for key := range flat {
			settingKeys[key] = struct{}{}
		}
	}
	sortedKeys := make([]string, 0, len(settingKeys))
	for key := range settingKeys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		columnDefs = append(columnDefs, output.ColumnDefaults{Header: key, Type: output.Text})
	}

	rows := make([]map[string]any, 0, len(repos))
	data := make([][]string, 0, len(repos))
	for name, repo := range repos {
		flat := flatSettings[name]
		row := map[string]any{
			"name": name,
			"type": repo.Type,
		}
		rowCells := make([]string, 0, len(columnDefs))
		rowCells = append(rowCells, name, repo.Type)
		for _, key := range sortedKeys {
			val, ok := flat[key]
			if ok {
				row[key] = val
				rowCells = append(rowCells, utils.FormatSettingValue(val))
			} else {
				rowCells = append(rowCells, "")
			}
		}
		rows = append(rows, row)
		data = append(data, rowCells)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(rows)
	}

	sortCols := output.ParseSortColumns("name")
	return output.PrintTable(columnDefs, data, sortCols)
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
