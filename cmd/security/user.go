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

func HandleUserList(ctx context.Context, nameFilter string) error {
	users, err := security.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}
	if nameFilter != "" {
		filtered := make(security.UserResponse)
		for name, user := range users {
			if strings.Contains(name, nameFilter) {
				filtered[name] = user
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("no users matched: %s", nameFilter)
		}
		return output.Render(filtered)
	}
	return output.Render(users)
}

func HandleUserGet(ctx context.Context, username string) error {
	user, err := security.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	return output.Render(user)
}

func HandleUserCreate(ctx context.Context, username, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var user security.User
	if err := json.Unmarshal(data, &user); err != nil {
		return fmt.Errorf("invalid user JSON: %w", err)
	}

	if err := security.CreateUser(ctx, username, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("User '%s' created/updated successfully\n", username)
	return nil
}

func HandleUserDelete(ctx context.Context, username string) error {
	if err := security.DeleteUser(ctx, username); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	fmt.Printf("User '%s' deleted successfully\n", username)
	return nil
}
