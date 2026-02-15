package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// indexNameRegex validates Elasticsearch index names
	// Must start with lowercase letter or digit, contain only lowercase letters, digits, hyphens, underscores
	indexNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

	// templateNameRegex is similar but also allows dots
	templateNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

	// aliasNameRegex for alias validation
	aliasNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

	// dataStreamStartRegex validates data stream names start with lowercase letter
	dataStreamStartRegex = regexp.MustCompile(`^[a-z]`)

	// dataStreamNameRegex validates full data stream name format
	dataStreamNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

	// timeUnitRegex validates time unit suffixes (for refresh intervals and timeouts)
	timeUnitRegex = regexp.MustCompile(`^.*?(ms|s|m|h|d)$`)
)

// ValidateIndexName validates an Elasticsearch index name
func ValidateIndexName(name string) error {
	if name == "" {
		return fmt.Errorf("index name cannot be empty")
	}

	// Check for reserved names
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("index name cannot start with '.' (reserved for system indices)")
	}

	if strings.HasPrefix(name, "_") {
		return fmt.Errorf("index name cannot start with '_' (reserved)")
	}

	// Check length
	if len(name) > 255 {
		return fmt.Errorf("index name too long: max 255 characters, got %d", len(name))
	}

	// Check format
	if !indexNameRegex.MatchString(name) {
		return fmt.Errorf("invalid index name: %s (must start with lowercase letter/digit, contain only lowercase letters, digits, hyphens, underscores)", name)
	}

	// Check for invalid characters (single pass)
	if strings.ContainsAny(name, `\/*?"<>| ,#`) {
		return fmt.Errorf("index name contains invalid characters")
	}

	return nil
}

// ValidateTemplateName validates an index template name
func ValidateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if len(name) > 255 {
		return fmt.Errorf("template name too long: max 255 characters, got %d", len(name))
	}

	if !templateNameRegex.MatchString(name) {
		return fmt.Errorf("invalid template name: %s", name)
	}

	return nil
}

// ValidateAliasName validates an alias name
func ValidateAliasName(name string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	if len(name) > 255 {
		return fmt.Errorf("alias name too long: max 255 characters, got %d", len(name))
	}

	if !aliasNameRegex.MatchString(name) {
		return fmt.Errorf("invalid alias name: %s", name)
	}

	return nil
}

// ValidateAliasPattern validates an alias pattern (with wildcards).
func ValidateAliasPattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("alias pattern cannot be empty")
	}

	if len(pattern) > 255 {
		return fmt.Errorf("alias pattern too long: max 255 characters, got %d", len(pattern))
	}

	// Patterns can contain * and ? wildcards.
	cleaned := strings.ReplaceAll(pattern, "*", "a")
	cleaned = strings.ReplaceAll(cleaned, "?", "a")

	if !aliasNameRegex.MatchString(cleaned) {
		return fmt.Errorf("invalid alias pattern: %s", pattern)
	}

	return nil
}

// ValidateShardCount validates the number of shards
func ValidateShardCount(count int) error {
	if count < 1 {
		return fmt.Errorf("shard count must be at least 1, got %d", count)
	}

	if count > 1024 {
		return fmt.Errorf("shard count too high: max 1024, got %d", count)
	}

	return nil
}

// ValidateReplicaCount validates the number of replicas
func ValidateReplicaCount(count int) error {
	if count < 0 {
		return fmt.Errorf("replica count cannot be negative, got %d", count)
	}

	if count > 100 {
		return fmt.Errorf("replica count too high: max 100, got %d", count)
	}

	return nil
}

// ValidateRefreshInterval validates refresh interval
func ValidateRefreshInterval(interval string) error {
	if interval == "" {
		return fmt.Errorf("refresh interval cannot be empty")
	}

	// Allow -1 (disable) or valid time units
	if interval == "-1" {
		return nil
	}

	// Check for valid time unit suffixes using regex
	if !timeUnitRegex.MatchString(interval) {
		return fmt.Errorf("invalid refresh interval: %s (must end with ms, s, m, h, or d, or be -1)", interval)
	}

	return nil
}

// ValidateTimeout validates timeout duration string
func ValidateTimeout(timeout string) error {
	if timeout == "" {
		return fmt.Errorf("timeout cannot be empty")
	}

	// Timeout uses same pattern as refresh interval (excluding days)
	if !strings.HasSuffix(timeout, "ms") && !strings.HasSuffix(timeout, "s") &&
		!strings.HasSuffix(timeout, "m") && !strings.HasSuffix(timeout, "h") {
		return fmt.Errorf("invalid timeout: %s (must end with ms, s, m, or h)", timeout)
	}

	return nil
}

// ValidateHostPort validates host:port format
func ValidateHostPort(hostPort string) error {
	if hostPort == "" {
		return fmt.Errorf("host:port cannot be empty")
	}

	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid host:port format: %s (expected host:port)", hostPort)
	}

	host := parts[0]
	if host == "" {
		return fmt.Errorf("host cannot be empty in %s", hostPort)
	}

	port := parts[1]
	if port == "" {
		return fmt.Errorf("port cannot be empty in %s", hostPort)
	}

	return nil
}

// ValidateURL validates a basic URL format
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid URL: %s (must start with http:// or https://)", url)
	}

	return nil
}

// ValidatePriority validates index template priority
func ValidatePriority(priority int) error {
	if priority < 0 {
		return fmt.Errorf("priority cannot be negative, got %d", priority)
	}

	if priority > 1000000 {
		return fmt.Errorf("priority too high: max 1000000, got %d", priority)
	}

	return nil
}

// ValidateIndexPattern validates index pattern (with wildcards)
func ValidateIndexPattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("index pattern cannot be empty")
	}

	if len(pattern) > 255 {
		return fmt.Errorf("index pattern too long: max 255 characters, got %d", len(pattern))
	}

	// Patterns can contain * and ? wildcards
	// Remove wildcards for validation
	cleaned := strings.ReplaceAll(pattern, "*", "a")
	cleaned = strings.ReplaceAll(cleaned, "?", "a")

	// Check if the cleaned pattern would be valid
	if !indexNameRegex.MatchString(cleaned) {
		return fmt.Errorf("invalid index pattern: %s", pattern)
	}

	return nil
}

// ValidateDataStreamName validates a data stream name
func ValidateDataStreamName(name string) error {
	if name == "" {
		return fmt.Errorf("data stream name cannot be empty")
	}

	// Check length
	if len(name) > 255 {
		return fmt.Errorf("data stream name too long: max 255 characters, got %d", len(name))
	}

	// Data streams must start with a lowercase letter
	if !dataStreamStartRegex.MatchString(name) {
		return fmt.Errorf("data stream name must start with a lowercase letter: %s", name)
	}

	// Data streams can contain lowercase letters, digits, hyphens, and underscores
	if !dataStreamNameRegex.MatchString(name) {
		return fmt.Errorf("invalid data stream name: %s (must contain only lowercase letters, digits, hyphens, underscores)", name)
	}

	// Check for invalid characters (single pass)
	if strings.ContainsAny(name, `\/*?"<>| ,#:`) {
		return fmt.Errorf("data stream name contains invalid characters")
	}

	return nil
}

// ValidateScriptLanguage validates a script language
func ValidateScriptLanguage(lang string) error {
	if lang == "" {
		return fmt.Errorf("script language cannot be empty")
	}

	// Valid script languages in Elasticsearch/OpenSearch
	validLangs := []string{
		"painless",   // Default language
		"mustache",   // Template language
		"expression", // Simple expressions
		"java",       // Java language (if enabled)
	}

	for _, valid := range validLangs {
		if strings.EqualFold(lang, valid) {
			return nil
		}
	}

	return fmt.Errorf("invalid script language: %s (valid: painless, mustache, expression, java)", lang)
}

// ValidateComponentTemplateName validates a component template name
func ValidateComponentTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("component template name cannot be empty")
	}

	// Component templates follow the same naming rules as index templates
	if len(name) > 255 {
		return fmt.Errorf("component template name too long: max 255 characters, got %d", len(name))
	}

	// Must start with lowercase letter or digit
	if !templateNameRegex.MatchString(name) {
		return fmt.Errorf("invalid component template name: %s (must start with lowercase letter/digit, contain only lowercase letters, digits, dots, hyphens, underscores)", name)
	}

	return nil
}

// ValidateThreadType validates a thread type for hot threads API
func ValidateThreadType(threadType string) error {
	if threadType == "" {
		// Empty is valid - uses default (cpu)
		return nil
	}

	validTypes := []string{"cpu", "wait", "block"}

	for _, valid := range validTypes {
		if strings.EqualFold(threadType, valid) {
			return nil
		}
	}

	return fmt.Errorf("invalid thread type: %s (valid: cpu, wait, block)", threadType)
}

// ValidateShardsForSplit validates shard counts for split operation
func ValidateShardsForSplit(sourceShard int, targetShard int) error {
	if sourceShard < 1 {
		return fmt.Errorf("source shard count must be at least 1, got %d", sourceShard)
	}

	if targetShard < 1 {
		return fmt.Errorf("target shard count must be at least 1, got %d", targetShard)
	}

	// Target must be a multiple of source
	if targetShard%sourceShard != 0 {
		return fmt.Errorf("target shard count (%d) must be a multiple of source shard count (%d)", targetShard, sourceShard)
	}

	// Target must be greater than source
	if targetShard <= sourceShard {
		return fmt.Errorf("target shard count (%d) must be greater than source shard count (%d)", targetShard, sourceShard)
	}

	return nil
}

// ValidateShardsForShrink validates shard counts for shrink operation
func ValidateShardsForShrink(sourceShard int, targetShard int) error {
	if sourceShard < 1 {
		return fmt.Errorf("source shard count must be at least 1, got %d", sourceShard)
	}

	if targetShard < 1 {
		return fmt.Errorf("target shard count must be at least 1, got %d", targetShard)
	}

	// Source must be a multiple of target
	if sourceShard%targetShard != 0 {
		return fmt.Errorf("source shard count (%d) must be a multiple of target shard count (%d)", sourceShard, targetShard)
	}

	// Target must be less than source
	if targetShard >= sourceShard {
		return fmt.Errorf("target shard count (%d) must be less than source shard count (%d)", targetShard, sourceShard)
	}

	return nil
}
