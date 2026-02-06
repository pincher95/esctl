package validation

import (
	"testing"
)

func TestValidateIndexName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "myindex", false},
		{"valid with hyphen", "my-index", false},
		{"valid with underscore", "my_index", false},
		{"valid with number", "index123", false},
		{"empty", "", true},
		{"starts with dot", ".system", true},
		{"starts with underscore", "_internal", true},
		{"too long", string(make([]byte, 256)), true},
		{"uppercase", "MyIndex", true},
		{"space", "my index", true},
		{"slash", "my/index", true},
		{"asterisk", "my*index", true},
		{"question mark", "my?index", true},
		{"comma", "my,index", true},
		{"hash", "my#index", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIndexName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIndexName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "mytemplate", false},
		{"valid with dot", "my.template", false},
		{"valid with hyphen", "my-template", false},
		{"empty", "", true},
		{"too long", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateAliasName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "myalias", false},
		{"valid with dot", "my.alias", false},
		{"valid with hyphen", "my-alias", false},
		{"empty", "", true},
		{"too long", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAliasName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAliasName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateAliasPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "myalias", false},
		{"valid with wildcard", "alias-*", false},
		{"valid with question mark", "alias-?", false},
		{"valid with dot", "my.alias*", false},
		{"empty", "", true},
		{"too long", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAliasPattern(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAliasPattern(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateShardCount(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"valid 1", 1, false},
		{"valid 5", 5, false},
		{"valid 100", 100, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"too high", 1025, true},
		{"max valid", 1024, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateShardCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateShardCount(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateReplicaCount(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"valid 0", 0, false},
		{"valid 1", 1, false},
		{"valid 5", 5, false},
		{"negative", -1, true},
		{"too high", 101, true},
		{"max valid", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateReplicaCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateReplicaCount(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRefreshInterval(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid seconds", "30s", false},
		{"valid milliseconds", "500ms", false},
		{"valid minutes", "5m", false},
		{"valid hours", "1h", false},
		{"valid days", "7d", false},
		{"disabled", "-1", false},
		{"empty", "", true},
		{"no unit", "30", true},
		{"invalid unit", "30x", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRefreshInterval(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRefreshInterval(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTimeout(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid seconds", "30s", false},
		{"valid milliseconds", "500ms", false},
		{"valid minutes", "5m", false},
		{"valid hours", "1h", false},
		{"empty", "", true},
		{"no unit", "30", true},
		{"invalid unit", "30d", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeout(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeout(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHostPort(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "localhost:9200", false},
		{"valid IP", "127.0.0.1:9200", false},
		{"valid domain", "es.example.com:9200", false},
		{"empty", "", true},
		{"no port", "localhost", true},
		{"no host", ":9200", true},
		{"multiple colons", "host:9200:extra", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHostPort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHostPort(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid http", "http://localhost:9200", false},
		{"valid https", "https://localhost:9200", false},
		{"empty", "", true},
		{"no protocol", "localhost:9200", true},
		{"invalid protocol", "ftp://localhost:9200", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePriority(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"valid 0", 0, false},
		{"valid 100", 100, false},
		{"valid 1000", 1000, false},
		{"negative", -1, true},
		{"too high", 1000001, true},
		{"max valid", 1000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePriority(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePriority(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateIndexPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "logs-*", false},
		{"valid multiple wildcards", "logs-*-*", false},
		{"valid question mark", "log?", false},
		{"valid mixed", "logs-*-202?", false},
		{"empty", "", true},
		{"too long", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIndexPattern(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIndexPattern(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
