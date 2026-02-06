package security

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type User struct {
	Username     string         `json:"username,omitempty"`
	Roles        []string       `json:"roles"`
	FullName     string         `json:"full_name,omitempty"`
	Email        string         `json:"email,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Enabled      bool           `json:"enabled"`
	Password     string         `json:"password,omitempty"`
	PasswordHash string         `json:"password_hash,omitempty"`
}

type UserResponse map[string]User

type Role struct {
	Cluster      []string          `json:"cluster,omitempty"`
	Indices      []IndexPrivilege  `json:"indices,omitempty"`
	Applications []ApplicationPriv `json:"applications,omitempty"`
	RunAs        []string          `json:"run_as,omitempty"`
	Metadata     map[string]any    `json:"metadata,omitempty"`
}

type IndexPrivilege struct {
	Names                  []string  `json:"names"`
	Privileges             []string  `json:"privileges"`
	FieldSecurity          *FieldSec `json:"field_security,omitempty"`
	Query                  string    `json:"query,omitempty"`
	AllowRestrictedIndices bool      `json:"allow_restricted_indices,omitempty"`
}

type FieldSec struct {
	Grant  []string `json:"grant,omitempty"`
	Except []string `json:"except,omitempty"`
}

type ApplicationPriv struct {
	Application string   `json:"application"`
	Privileges  []string `json:"privileges"`
	Resources   []string `json:"resources"`
}

type RoleResponse map[string]Role

// ListUsers lists all users
func ListUsers(ctx context.Context) (UserResponse, error) {
	var result UserResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_security/user")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list users: %s", resp.Status())
	}
	return result, nil
}

// GetUser gets a specific user
func GetUser(ctx context.Context, username string) (UserResponse, error) {
	var result UserResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_security/user/%s", username))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get user: %s", resp.Status())
	}
	return result, nil
}

// CreateUser creates or updates a user
func CreateUser(ctx context.Context, username string, user User) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(user).
		Put(fmt.Sprintf("_security/user/%s", username))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to create user: %s", resp.Status())
	}
	return nil
}

// DeleteUser deletes a user
func DeleteUser(ctx context.Context, username string) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("_security/user/%s", username))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete user: %s", resp.Status())
	}
	return nil
}

// ListRoles lists all roles
func ListRoles(ctx context.Context) (RoleResponse, error) {
	var result RoleResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_security/role")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list roles: %s", resp.Status())
	}
	return result, nil
}

// GetRole gets a specific role
func GetRole(ctx context.Context, name string) (RoleResponse, error) {
	var result RoleResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_security/role/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get role: %s", resp.Status())
	}
	return result, nil
}

// CreateRole creates or updates a role
func CreateRole(ctx context.Context, name string, role Role) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(role).
		Put(fmt.Sprintf("_security/role/%s", name))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to create role: %s", resp.Status())
	}
	return nil
}

// DeleteRole deletes a role
func DeleteRole(ctx context.Context, name string) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(fmt.Sprintf("_security/role/%s", name))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete role: %s", resp.Status())
	}
	return nil
}
