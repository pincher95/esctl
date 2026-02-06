package es

import "testing"

func TestGetNestedPath(t *testing.T) {
	tests := []struct {
		field       string
		nestedPaths []string
		wantPath    string
		wantOK      bool
	}{
		{"author.name", []string{"author"}, "author", true},
		{"author.address.city", []string{"author.address"}, "author.address", true},
		{"title", []string{"author"}, "", false},
		{"nested.field", []string{"other", "nested"}, "nested", true},
	}

	for _, tt := range tests {
		got, ok := getNestedPath(tt.field, tt.nestedPaths)
		if got != tt.wantPath || ok != tt.wantOK {
			t.Fatalf("getNestedPath(%q, %v) = (%q, %v), want (%q, %v)", tt.field, tt.nestedPaths, got, ok, tt.wantPath, tt.wantOK)
		}
	}
}

func TestMax(t *testing.T) {
	if max(1, 2) != 2 {
		t.Fatalf("max(1,2) != 2")
	}
	if max(5, 3) != 5 {
		t.Fatalf("max(5,3) != 5")
	}
	if max(7, 7) != 7 {
		t.Fatalf("max(7,7) != 7")
	}
}
