package security

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pincher95/esctl/es/security"
	"github.com/pincher95/esctl/output"
)

func HandleRoleList(ctx context.Context) error {
	roles, err := security.ListRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to list roles: %w", err)
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
