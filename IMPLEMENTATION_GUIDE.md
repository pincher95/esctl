# Implementation Guide for New Features

This guide shows you how to add new features to esctl following the established patterns.

## Table of Contents
- [Adding a New Command](#adding-a-new-command)
- [Adding ES/OS API Support](#adding-esos-api-support)
- [Writing Tests](#writing-tests)
- [Best Practices](#best-practices)

---

## Adding a New Command

### Step 1: Create the ES Package

Create the business logic in `es/<feature>/`:

```go
// es/ilm/ilm.go
package ilm

import (
	"context"
	"fmt"
	"github.com/pincher95/esctl/shared"
)

type Policy struct {
	Name   string                 `json:"-"`
	Policy map[string]interface{} `json:"policy"`
}

func List(ctx context.Context) (map[string]Policy, error) {
	var result map[string]Policy

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get("_ilm/policy")

	if err != nil {
		return nil, fmt.Errorf("failed to list ILM policies: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list ILM policies: %s", resp.Status())
	}

	return result, nil
}

func Get(ctx context.Context, name string) (*Policy, error) {
	var result map[string]Policy

	endpoint := fmt.Sprintf("_ilm/policy/%s", name)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to get ILM policy: %w", err)
	}

	if resp.StatusCode() == 404 {
		return nil, fmt.Errorf("ILM policy not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get ILM policy: %s", resp.Status())
	}

	if policy, ok := result[name]; ok {
		policy.Name = name
		return &policy, nil
	}

	return nil, fmt.Errorf("failed to parse ILM policy response")
}

func Put(ctx context.Context, name string, policy Policy) error {
	endpoint := fmt.Sprintf("_ilm/policy/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(policy).
		Put(endpoint)

	if err != nil {
		return fmt.Errorf("failed to put ILM policy: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to put ILM policy: %s - %s",
			resp.Status(), string(resp.Body()))
	}

	return nil
}

func Delete(ctx context.Context, name string) error {
	endpoint := fmt.Sprintf("_ilm/policy/%s", name)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		return fmt.Errorf("failed to delete ILM policy: %w", err)
	}

	if resp.StatusCode() == 404 {
		return fmt.Errorf("ILM policy not found: %s", name)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete ILM policy: %s", resp.Status())
	}

	return nil
}
```

### Step 2: Create Tests

Always create tests alongside your implementation:

```go
// es/ilm/ilm_test.go
package ilm

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestList(t *testing.T) {
	mockResp := map[string]interface{}{
		"hot_delete_policy": map[string]interface{}{
			"policy": map[string]interface{}{
				"phases": map[string]interface{}{
					"hot": map[string]interface{}{
						"actions": map[string]interface{}{
							"rollover": map[string]interface{}{
								"max_age": "7d",
							},
						},
					},
				},
			},
		},
	}

	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/_ilm/policy")
	defer srv.Close()
	shared.SetClient(cli)

	policies, err := List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(policies))
	}

	if _, ok := policies["hot_delete_policy"]; !ok {
		t.Error("Expected hot_delete_policy to be present")
	}
}

func TestGetNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404,
		`{"error":"not found"}`, "/_ilm/policy/nonexistent")
	defer srv.Close()
	shared.SetClient(cli)

	_, err := Get(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent policy")
	}
}
```

### Step 3: Create Command Structure

Create the command in `cmd/<feature>/`:

```go
// cmd/ilm/ilm.go
package ilm

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "ilm",
	Short: "Manage Index Lifecycle Management policies",
	Long: `Manage ILM policies in Elasticsearch.

ILM policies define how indices are managed over time, including
rollover, shrink, delete, and other lifecycle actions.`,
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(putCmd)
	Cmd.AddCommand(deleteCmd)
	Cmd.AddCommand(explainCmd)
}
```

```go
// cmd/ilm/list.go
package ilm

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all ILM policies",
	Example: utils.TrimAndIndent(`
		# List all ILM policies
		esctl ilm list

		# Output as JSON
		esctl ilm list -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		policies, err := ilm.List(ctx)
		if err != nil {
			return err
		}

		if len(policies) == 0 {
			fmt.Println("No ILM policies found")
			return nil
		}

		// For JSON/YAML output
		if output.Format != "table" {
			output.PrintJson(policies)
			return nil
		}

		// For table output
		columnDefs := []output.ColumnDefaults{
			{Header: "NAME", Type: output.Text},
			{Header: "PHASES", Type: output.Text},
		}

		var data [][]string
		for name := range policies {
			// You can extract more details from the policy
			data = append(data, []string{
				name,
				"hot,delete", // Parse from actual policy
			})
		}

		output.PrintTable(columnDefs, data, nil)
		return nil
	},
}
```

```go
// cmd/ilm/get.go
package ilm

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get details of a specific ILM policy",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Get ILM policy details
		esctl ilm get hot_delete_policy

		# Get as JSON
		esctl ilm get hot_delete_policy -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		policy, err := ilm.Get(ctx, name)
		if err != nil {
			return err
		}

		output.PrintJson(policy)
		return nil
	},
}
```

### Step 4: Register the Command

Add to `cmd/root.go`:

```go
import (
	// ... other imports
	"github.com/pincher95/esctl/cmd/ilm"
)

func init() {
	// ... other commands
	RootCmd.AddCommand(ilm.Cmd)
}
```

### Step 5: Test Manually

```bash
# Build
go build -o esctl

# Test
./esctl ilm list
./esctl ilm get my-policy
./esctl ilm put my-policy --file policy.json
./esctl ilm delete my-policy
```

---

## Adding ES/OS API Support

### HTTP Client Patterns

**GET Request:**
```go
var result YourStruct

resp, err := shared.Client.R().
	SetContext(ctx).
	SetHeader("Content-Type", "application/json").
	SetResult(&result).
	Get(endpoint)

if err != nil {
	return nil, fmt.Errorf("request failed: %w", err)
}

if resp.StatusCode() != 200 {
	return nil, fmt.Errorf("unexpected status: %s", resp.Status())
}
```

**POST Request:**
```go
var result YourStruct
body := map[string]interface{}{
	"key": "value",
}

resp, err := shared.Client.R().
	SetContext(ctx).
	SetHeader("Content-Type", "application/json").
	SetBody(body).
	SetResult(&result).
	Post(endpoint)
```

**PUT Request:**
```go
resp, err := shared.Client.R().
	SetContext(ctx).
	SetHeader("Content-Type", "application/json").
	SetBody(data).
	Put(endpoint)
```

**DELETE Request:**
```go
resp, err := shared.Client.R().
	SetContext(ctx).
	SetHeader("Content-Type", "application/json").
	Delete(endpoint)
```

**HEAD Request (exists check):**
```go
resp, err := shared.Client.R().
	SetContext(ctx).
	Head(endpoint)

return resp.StatusCode() == 200, err
```

### Query Parameters

```go
import "net/url"

u := url.URL{Path: "_search"}
q := u.Query()
q.Set("size", "100")
q.Set("from", "0")
u.RawQuery = q.Encode()

endpoint := u.String() // "_search?from=0&size=100"
```

---

## Writing Tests

### Basic Test Structure

```go
package myfeature

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestMyFeature(t *testing.T) {
	// 1. Create mock response
	mockResp := map[string]interface{}{
		"acknowledged": true,
	}

	// 2. Setup mock server
	respJSON, _ := json.Marshal(mockResp)
	srv, cli := testutil.NewMockServer(string(respJSON), "/expected/path")
	defer srv.Close()

	// 3. Inject mock client
	shared.SetClient(cli)

	// 4. Call function
	err := MyFeature(context.Background())

	// 5. Assert results
	if err != nil {
		t.Fatalf("MyFeature() error = %v", err)
	}
}
```

### Testing Error Cases

```go
func TestMyFeatureNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404,
		`{"error":"not found"}`, "/expected/path")
	defer srv.Close()
	shared.SetClient(cli)

	_, err := MyFeature(context.Background())
	if err == nil {
		t.Error("Expected error for 404 response")
	}
}
```

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./es/ilm

# Verbose
go test -v ./es/ilm

# With coverage
go test -cover ./...

# Single test
go test ./es/ilm -run TestList -v
```

---

## Best Practices

### 1. Error Handling

**Always wrap errors with context:**
```go
if err != nil {
	return fmt.Errorf("failed to do something: %w", err)
}
```

**Check HTTP status codes:**
```go
if resp.StatusCode() == 404 {
	return fmt.Errorf("resource not found: %s", name)
}

if resp.StatusCode() != 200 {
	return fmt.Errorf("unexpected status: %s - %s",
		resp.Status(), string(resp.Body()))
}
```

### 2. Context Propagation

**Always accept context as first parameter:**
```go
func MyFunction(ctx context.Context, name string) error {
	// Use ctx in all HTTP calls
	resp, err := shared.Client.R().
		SetContext(ctx).
		Get(endpoint)
	// ...
}
```

### 3. Output Formatting

**Support multiple output formats:**
```go
// For detailed output (single resource)
output.PrintJson(resource)

// For lists with table support
if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
	if shared.OutputFormat == "json" {
		output.PrintJson(resources)
	} else {
		output.PrintYaml(resources)
	}
	return nil
}

// Table format
columnDefs := []output.ColumnDefaults{
	{Header: "NAME", Type: output.Text},
	{Header: "VALUE", Type: output.Number},
}
output.PrintTable(columnDefs, data, nil)
```

### 4. Command Structure

**Use clear, descriptive names:**
- `list` - List all resources
- `get <name>` - Get single resource
- `put <name>` - Create/update resource
- `delete <name>` - Delete resource
- `exists <name>` - Check existence

**Add examples:**
```go
Example: utils.TrimAndIndent(`
	# Basic usage
	esctl command subcommand

	# With flags
	esctl command subcommand --flag value

	# Complex example
	esctl command subcommand --multiple --flags
`),
```

### 5. Flag Naming

**Use consistent flag names:**
- `--index` for index names
- `--node` for node names
- `--file` for file input
- `--output` or `-o` for output format
- `--force` to skip confirmations
- `--timeout` for operation timeout

### 6. Confirmation for Destructive Operations

**Always confirm unless --force:**
```go
if !flagForce {
	fmt.Printf("This will delete %s. Continue?\n", name)
	approved, err := utils.GetApproval()
	if err != nil || !approved {
		fmt.Println("Operation cancelled")
		return nil
	}
}
```

### 7. JSON Structs

**Use struct tags properly:**
```go
type MyStruct struct {
	Name      string                 `json:"name"`
	Count     int                    `json:"count,omitempty"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
	Internal  string                 `json:"-"` // Never marshaled
}
```

### 8. Package Organization

```
es/
└── myfeature/
    ├── myfeature.go       # Main implementation
    ├── myfeature_test.go  # Tests
    └── types.go           # Type definitions (if many)

cmd/
└── myfeature/
    ├── myfeature.go       # Main command
    ├── list.go            # List subcommand
    ├── get.go             # Get subcommand
    ├── put.go             # Put subcommand
    └── delete.go          # Delete subcommand
```

---

## Reference Implementation

See the **template** command for a complete reference implementation:

- `es/template/template.go` - ES API layer
- `es/template/template_test.go` - Tests
- `cmd/template/*.go` - CLI commands

This implementation demonstrates:
- ✅ Proper error handling
- ✅ Context usage
- ✅ Multiple output formats
- ✅ Comprehensive tests
- ✅ User confirmations
- ✅ Good documentation

---

## Checklist for New Features

Before submitting:

- [ ] ES package implementation in `es/<feature>/`
- [ ] Comprehensive tests with >80% coverage
- [ ] CLI commands in `cmd/<feature>/`
- [ ] Command registered in `cmd/root.go`
- [ ] Examples in command help text
- [ ] Error handling with context
- [ ] Context propagation
- [ ] Support for table/json/yaml output
- [ ] Confirmation for destructive operations
- [ ] Linter passes (`go vet ./...`)
- [ ] All tests pass (`go test ./...`)
- [ ] Manual testing completed

---

## Getting Help

- Check existing commands for patterns
- Review `es/template/` for a complete example
- Look at test files for testing patterns
- See `REVIEW_AND_RECOMMENDATIONS.md` for architecture guidance

Happy coding! 🚀
