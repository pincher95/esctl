# esctl Code Review & Recommendations

## Executive Summary

**esctl** is a well-structured CLI tool for managing Elasticsearch/OpenSearch clusters. The codebase demonstrates good practices in separation of concerns, testing, and user experience. This document outlines suggested improvements and new features to enhance the tool's capabilities.

---

## 📊 Architecture Assessment

### Strengths
- ✅ Clean separation: `cmd/` (CLI), `es/` (business logic), `output/` (formatting)
- ✅ Consistent use of Cobra for CLI commands
- ✅ Context propagation for timeout/cancellation
- ✅ Multiple output formats (table, JSON, YAML)
- ✅ Configuration management with contexts
- ✅ Good test coverage with mock HTTP servers
- ✅ Watch mode for live monitoring

### Areas for Improvement
- ⚠️ Inconsistent error handling across commands
- ⚠️ No structured logging framework
- ⚠️ Missing rate limiting for bulk operations
- ⚠️ No connection pooling configuration
- ⚠️ Limited input validation

---

## 🔧 Recommended Improvements

### 1. Enhanced Error Handling

**Status**: ✅ Created `internal/errors/errors.go`

Benefits:
- Structured error types for better debugging
- User-friendly error messages
- Retry logic for transient failures
- Authentication error detection

### 2. Logging Framework

Add structured logging for debugging and operations:

```go
// internal/logger/logger.go
package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}
```

### 3. Input Validation

Add validation helpers:

```go
// internal/validation/validation.go
package validation

import (
	"fmt"
	"regexp"
)

var indexNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

func ValidateIndexName(name string) error {
	if name == "" {
		return fmt.Errorf("index name cannot be empty")
	}
	if !indexNameRegex.MatchString(name) {
		return fmt.Errorf("invalid index name: %s (must start with lowercase letter/digit, contain only lowercase letters, digits, hyphens, underscores)", name)
	}
	if len(name) > 255 {
		return fmt.Errorf("index name too long: max 255 characters")
	}
	return nil
}

func ValidateShardCount(count int) error {
	if count < 1 {
		return fmt.Errorf("shard count must be at least 1")
	}
	if count > 1024 {
		return fmt.Errorf("shard count too high: max 1024")
	}
	return nil
}
```

### 4. Progress Indicators

For long-running operations (reindex, forcemerge, etc.):

```go
// internal/progress/progress.go
package progress

import (
	"fmt"
	"time"
)

type Bar struct {
	total   int64
	current int64
	start   time.Time
}

func New(total int64) *Bar {
	return &Bar{total: total, start: time.Now()}
}

func (b *Bar) Add(n int64) {
	b.current += n
	b.render()
}

func (b *Bar) render() {
	percent := float64(b.current) / float64(b.total) * 100
	elapsed := time.Since(b.start)
	fmt.Printf("\rProgress: %.1f%% (%d/%d) - Elapsed: %s",
		percent, b.current, b.total, elapsed.Round(time.Second))
}
```

### 5. Rate Limiting

For bulk operations to prevent cluster overload:

```go
// internal/ratelimit/limiter.go
package ratelimit

import (
	"context"
	"time"
)

type Limiter struct {
	ticker *time.Ticker
}

func New(requestsPerSecond int) *Limiter {
	interval := time.Second / time.Duration(requestsPerSecond)
	return &Limiter{ticker: time.NewTicker(interval)}
}

func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Stop() {
	l.ticker.Stop()
}
```

---

## 🆕 New Commands to Add

### Priority 1: Essential Operations

#### 1. Index Templates Management

```bash
# List templates
esctl template list

# Get template details
esctl template get <name>

# Create/update template
esctl template put <name> --file template.json

# Delete template
esctl template delete <name>
```

**Implementation:**
```go
// cmd/template/template.go
package template

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "template",
	Short: "Manage index templates",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(getCmd)
	Cmd.AddCommand(putCmd)
	Cmd.AddCommand(deleteCmd)
}
```

#### 2. Index Lifecycle Management (ILM)

```bash
# List ILM policies
esctl ilm list

# Get policy details
esctl ilm get <policy>

# Create/update policy
esctl ilm put <policy> --file policy.json

# Explain index lifecycle status
esctl ilm explain --index <index>

# Retry failed ILM steps
esctl ilm retry --index <index>
```

#### 3. Cluster Backup/Restore Enhancement

```bash
# Create repository
esctl set snapshot-repo <name> --type fs --settings "location:/backup"

# Restore with options
esctl update snapshot <repo> <snapshot> \
  --indices "index1,index2" \
  --rename-pattern "(.+)" \
  --rename-replacement "restored_$1" \
  --include-global-state=false
```

#### 4. Data Stream Management

```bash
# List data streams
esctl get datastreams

# Get data stream details
esctl get datastream <name>

# Create data stream
esctl set datastream <name> --template <template>

# Delete data stream
esctl delete datastream <name>

# Rollover data stream
esctl update datastream rollover <name>
```

#### 5. Index Statistics Deep Dive

```bash
# Detailed index stats
esctl stats index <index> --groups indexing,search,merge

# Segment details
esctl stats segments <index>

# Recovery status
esctl stats recovery <index>
```

### Priority 2: Advanced Features

#### 6. Security Management (OpenSearch/X-Pack)

```bash
# User management
esctl get users
esctl set user <username> --file user.json

# Role management
esctl get roles
esctl get role <role>
```

#### 7. Index Disk Usage Analysis

```bash
# Analyze disk usage by index/field
esctl analyze disk --index <index> --top 20

# Find large segments
esctl analyze segments --min-size 1GB
```

#### 8. Query Performance Analysis

```bash
# Show slow queries
esctl analyze slow-queries --index <index> --threshold 1s

# Hot threads analysis
esctl analyze hot-threads --threads 5
```

#### 9. Cluster Diagnostics

```bash
# Health check with recommendations
esctl diagnose cluster

# Check for common issues
esctl diagnose shards --unassigned

# Disk watermark warnings
esctl diagnose disk
```

#### 10. Index Comparison

```bash
# Compare mappings between indices
esctl compare mappings index1 index2

# Compare settings
esctl compare settings index1 index2

# Compare data (document counts, sizes)
esctl compare data index1 index2
```

### Priority 3: Developer Tools

#### 11. Index Migration Helper

```bash
# Plan migration from old to new index
esctl migrate plan old_index new_index

# Execute migration with progress
esctl migrate execute old_index new_index --alias shared_alias
```

#### 12. Bulk Import/Export

```bash
# Export to ndjson
esctl export --index articles --output articles.ndjson

# Import from ndjson
esctl import --index articles --file articles.ndjson --batch-size 1000
```

#### 13. Cluster Comparison

```bash
# Compare two clusters
esctl compare clusters --source prod --target staging

# Sync configurations
esctl sync templates --from prod --to staging --dry-run
```

---

## 🏗️ Restructuring Suggestions

### 1. Command Organization

**Current Structure:**
```
cmd/
├── get/          # Read operations
├── describe/     # Detailed views
├── update/       # Update operations
├── set/          # Set configurations
└── delete/       # Delete operations
```

**Suggested Structure (Resource-Based):**
```
cmd/
├── index/
│   ├── create.go
│   ├── delete.go
│   ├── list.go
│   ├── update.go
│   ├── flush.go
│   ├── refresh.go
│   └── forcemerge.go
├── cluster/
│   ├── health.go
│   ├── settings.go
│   ├── stats.go
│   └── reroute.go
├── template/
│   └── ... (new)
├── ilm/
│   └── ... (new)
└── datastream/
    └── ... (new)
```

**Benefits:**
- Easier to find related commands
- Better grouping by resource type
- Aligns with Elasticsearch API structure

### 2. ES Package Organization

**Current:**
```
es/
├── cat/          # CAT API calls
├── cluster/      # Cluster APIs
├── index/        # Index APIs
└── search.go     # Search APIs
```

**Suggested Enhancement:**
```
es/
├── client/       # HTTP client wrappers
├── cat/          # CAT APIs
├── cluster/      # Cluster APIs
├── index/        # Index APIs
├── template/     # Template APIs (new)
├── ilm/          # ILM APIs (new)
├── datastream/   # Data stream APIs (new)
├── security/     # Security APIs (new)
└── common/       # Shared types
```

### 3. Shared Utilities

Create domain-specific helper packages:

```
internal/
├── client/       # HTTP client
├── errors/       # Error types ✅ (created)
├── logger/       # Structured logging
├── progress/     # Progress bars
├── ratelimit/    # Rate limiting
├── validation/   # Input validation
└── testutil/     # Test helpers
```

### 4. Configuration Management

Enhance the config system:

```go
// shared/config.go - Enhanced
type ClusterConfig struct {
	Name             string
	Protocol         string
	Host             string
	Port             int
	Username         string
	Password         string
	APIKey           string        // Support API key auth
	Headers          map[string]string // Custom headers
	SkipTLSVerify    bool
	CertFile         string
	Timeout          time.Duration
	MaxRetries       int

	// Display preferences
	DefaultOutput    string
	ColumnPresets    map[string][]string
}
```

---

## 📈 Performance Optimizations

### 1. Connection Pooling

Configure proper connection pooling in the HTTP client:

```go
r.SetTransport(&http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     90 * time.Second,
})
```

### 2. Pagination

For large result sets:

```go
// cmd/get/indices.go
var (
	flagLimit  int
	flagOffset int
)

// Support pagination
getCmd.Flags().IntVar(&flagLimit, "limit", 100, "Maximum results")
getCmd.Flags().IntVar(&flagOffset, "offset", 0, "Results offset")
```

### 3. Caching

Cache cluster metadata for short periods:

```go
// internal/cache/cache.go
type Cache struct {
	mu    sync.RWMutex
	data  map[string]cacheEntry
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiration) {
		return nil, false
	}
	return entry.value, true
}
```

---

## 🧪 Testing Enhancements

### 1. Integration Tests

Create integration test suite:

```go
// tests/integration/integration_test.go
// +build integration

func TestIndexLifecycle(t *testing.T) {
	// Requires ESCTL_TEST_HOST env var
	// Full lifecycle: create, insert, query, delete
}
```

Run with: `go test -tags=integration ./tests/integration`

### 2. E2E Tests

Create end-to-end CLI tests:

```bash
#!/bin/bash
# tests/e2e/test_index_operations.sh

set -e

# Create index
esctl index create test_index

# Verify creation
esctl get indices | grep test_index

# Delete index
esctl delete index test_index
```

### 3. Benchmark Tests

Add performance benchmarks:

```go
func BenchmarkListIndices(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = cat.Indices(context.Background())
	}
}
```

---

## 📝 Documentation Enhancements

### 1. Command Examples

Add more examples to each command:

```go
var Cmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze text with the specified analyzer",
	Example: utils.TrimAndIndent(`
		# Standard analyzer
		esctl analyze --text="The quick brown fox"

		# With specific analyzer
		esctl analyze --analyzer=whitespace --text="hello world"

		# Using field analyzer
		esctl analyze --index=articles --field=title --text="Elasticsearch Tips"

		# Multiple texts
		esctl analyze --text="text1" --text="text2"
	`),
}
```

### 2. Man Pages

Generate man pages:

```go
// cmd/docs/docs.go
func generateManPages() error {
	header := &doc.GenManHeader{
		Title:   "ESCTL",
		Section: "1",
	}
	return doc.GenManTree(cmd.RootCmd, header, "./man")
}
```

### 3. Completion Scripts

Add shell completion generation:

```go
// cmd/completion/completion.go
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		return fmt.Errorf("unsupported shell: %s", args[0])
	},
}
```

---

## 🔐 Security Enhancements

### 1. Secure Credential Storage

Use system keyring for credentials:

```go
// Add dependency: github.com/zalando/go-keyring

func storeCredentials(context, username, password string) error {
	return keyring.Set("esctl", context, fmt.Sprintf("%s:%s", username, password))
}

func getCredentials(context string) (string, string, error) {
	cred, err := keyring.Get("esctl", context)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(cred, ":", 2)
	return parts[0], parts[1], nil
}
```

### 2. API Key Support

Add API key authentication:

```go
// shared/config.go
type Context struct {
	// ... existing fields
	APIKey    string
	APISecret string
}

// internal/client/resty.go
if cfg.APIKey != "" {
	r.SetHeader("Authorization", "ApiKey "+cfg.APIKey)
}
```

### 3. TLS Configuration

Add TLS options:

```go
type Config struct {
	// ... existing fields
	SkipTLSVerify bool
	CertFile      string
	KeyFile       string
	CAFile        string
}
```

---

## 🎯 UX Improvements

### 1. Interactive Mode

Add interactive prompts for destructive operations:

```go
// Already have utils.GetApproval(), use it more:
if !flagForce {
	approved, err := utils.GetApproval()
	if err != nil || !approved {
		return fmt.Errorf("operation cancelled")
	}
}
```

### 2. Colored Output

Add color support:

```go
// Use fatih/color
import "github.com/fatih/color"

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

// In health display:
switch status {
case "green":
	fmt.Printf("Status: %s\n", green(status))
case "yellow":
	fmt.Printf("Status: %s\n", yellow(status))
case "red":
	fmt.Printf("Status: %s\n", red(status))
}
```

### 3. Smart Defaults

Improve default behaviors:

```go
// Auto-detect output format based on terminal
if !terminal.IsTerminal(int(os.Stdout.Fd())) {
	// Piped output - use JSON
	shared.OutputFormat = "json"
}
```

---

## 🎪 Advanced Features

### 1. Plugin System

Allow custom commands:

```go
// internal/plugin/plugin.go
type Plugin interface {
	Name() string
	Command() *cobra.Command
}

func LoadPlugins(dir string) ([]Plugin, error) {
	// Load .so files or run external binaries
}
```

### 2. Scripting Support

Execute multiple commands:

```bash
# script.esctl
index create articles --shards 3 --replicas 1
bulk insert articles --file data.json
index refresh articles
get indices --index articles
```

```bash
esctl script run script.esctl
```

### 3. Watch with Alerts

Monitor and alert on conditions:

```bash
# Alert when unassigned shards appear
esctl watch shards --unassigned --alert-on-change \
  --notify-email ops@example.com
```

---

## 🔍 Monitoring & Observability

### 1. Metrics Export

Export metrics in Prometheus format:

```go
// cmd/metrics/metrics.go
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Export cluster metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Collect metrics
		// Output in Prometheus format
	},
}
```

### 2. Health Checks

Comprehensive health checks:

```bash
esctl health check --verbose
```

Output:
```
✓ Cluster status: green
✓ All nodes responding
✓ No unassigned shards
✓ Disk usage below 80%
⚠ Heap usage at 85% on node-3
✗ JVM old gen GC taking >1s on node-1
```

---

## 📊 Summary of Changes

### Completed
- ✅ Enhanced HTTP client with retry logic and auth
- ✅ Created structured error handling package
- ✅ Improved timeout and retry configuration

### Recommended Implementation Priority

**Phase 1 (Essential - 2-3 weeks)**
1. Add structured logging
2. Implement input validation
3. Add progress indicators
4. Create index template management commands
5. Add ILM policy management

**Phase 2 (Important - 3-4 weeks)**
6. Restructure commands (resource-based)
7. Add data stream support
8. Implement rate limiting
9. Enhanced error messages
10. Add completion scripts

**Phase 3 (Advanced - 4-6 weeks)**
11. Add security management commands
12. Implement caching layer
13. Add bulk import/export
14. Create diagnostic tools
15. Add integration tests

**Phase 4 (Polish - 2-3 weeks)**
16. Colored output
17. Better documentation
18. Man pages
19. Plugin system foundation
20. Performance benchmarks

---

## 🎓 Learning Resources

For contributors to understand Elasticsearch/OpenSearch better:

1. **Elasticsearch API Documentation**: https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html
2. **OpenSearch API Documentation**: https://opensearch.org/docs/latest/api-reference/
3. **Go CLI Best Practices**: https://github.com/spf13/cobra/blob/master/user_guide.md

---

## 🤝 Contributing Guidelines

When implementing new features:

1. **Follow existing patterns**: Look at similar commands for structure
2. **Add tests**: Every new feature needs unit tests
3. **Update docs**: Add examples to command help and README
4. **Error handling**: Use structured errors with helpful messages
5. **Context support**: Always accept and propagate context
6. **Output formats**: Support table, JSON, and YAML where applicable

---

## 📞 Questions?

This review is comprehensive but not exhaustive. For specific questions about implementation details or priorities, feel free to ask!
