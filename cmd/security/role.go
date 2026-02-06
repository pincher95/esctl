package security

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pincher95/esctl/es/security"
	"github.com/pincher95/esctl/output"
)

func HandleRoleList(ctx context.Context, nameFilter string) error {
	roles, err := security.ListRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to list roles: %w", err)
	}
	if nameFilter != "" {
		filtered := make(security.RoleResponse)
		for name, role := range roles {
			if strings.Contains(name, nameFilter) {
				filtered[name] = role
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("no roles matched: %s", nameFilter)
		}
		return output.Render(filtered)
	}
	return output.Render(roles)
}

func HandleRoleGet(ctx context.Context, name string) error {
	role, err := security.GetRole(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	return output.Render(role)
}

func HandleRoleCreate(ctx context.Context, name, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var role security.Role
	if err := json.Unmarshal(data, &role); err != nil {
		return fmt.Errorf("invalid role JSON: %w", err)
	}

	if err := security.CreateRole(ctx, name, role); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	fmt.Printf("Role '%s' created/updated successfully\n", name)
	return nil
}

func HandleRoleDelete(ctx context.Context, name string) error {
	if err := security.DeleteRole(ctx, name); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	fmt.Printf("Role '%s' deleted successfully\n", name)
	return nil
}
