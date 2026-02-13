package script

import (
	"context"
	"testing"

	"github.com/pincher95/esctl/internal/testutil"
	"github.com/pincher95/esctl/shared"
)

func TestList(t *testing.T) {
	mockResponse := `{
		"script1": {
			"found": true,
			"_id": "script1",
			"script": {
				"lang": "painless",
				"source": "Math.log(_score * 2) + params.multiplier"
			}
		},
		"script2": {
			"found": true,
			"_id": "script2",
			"script": {
				"lang": "mustache",
				"source": "{\"query\":{\"match\":{\"{{field}}\":\"{{value}}\"}}}"
			}
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts")
	defer srv.Close()

	shared.SetClient(cli)

	scripts, err := List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(scripts) != 2 {
		t.Fatalf("expected 2 scripts, got %d", len(scripts))
	}

	if scripts["script1"].Script.Lang != "painless" {
		t.Errorf("expected script1 lang to be painless, got %s", scripts["script1"].Script.Lang)
	}

	if scripts["script2"].Script.Lang != "mustache" {
		t.Errorf("expected script2 lang to be mustache, got %s", scripts["script2"].Script.Lang)
	}
}

func TestGet(t *testing.T) {
	mockResponse := `{
		"found": true,
		"_id": "my-script",
		"script": {
			"lang": "painless",
			"source": "Math.log(_score * 2) + params.multiplier"
		}
	}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts/my-script")
	defer srv.Close()

	shared.SetClient(cli)

	script, err := Get(context.Background(), "my-script")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if script.Script.Lang != "painless" {
		t.Errorf("expected lang to be painless, got %s", script.Script.Lang)
	}

	if script.Script.Source != "Math.log(_score * 2) + params.multiplier" {
		t.Errorf("unexpected source: %s", script.Script.Source)
	}

	if script.ID != "my-script" {
		t.Errorf("expected ID to be my-script, got %s", script.ID)
	}
}

func TestGetNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_scripts/nonexistent")
	defer srv.Close()

	shared.SetClient(cli)

	_, err := Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent script")
	}

	if err.Error() != "script 'nonexistent' not found" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetEmptyID(t *testing.T) {
	_, err := Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}

	if err.Error() != "script ID is required" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPut(t *testing.T) {
	mockResponse := `{"acknowledged": true}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts/my-script")
	defer srv.Close()

	shared.SetClient(cli)

	script := Script{
		Lang:   "painless",
		Source: "Math.log(_score * 2) + params.multiplier",
	}

	err := Put(context.Background(), "my-script", script)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPutValidation(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		script      Script
		expectedErr string
	}{
		{
			name:        "empty ID",
			id:          "",
			script:      Script{Lang: "painless", Source: "test"},
			expectedErr: "script ID is required",
		},
		{
			name:        "empty language",
			id:          "test",
			script:      Script{Lang: "", Source: "test"},
			expectedErr: "script language is required",
		},
		{
			name:        "empty source",
			id:          "test",
			script:      Script{Lang: "painless", Source: ""},
			expectedErr: "script source is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Put(context.Background(), tt.id, tt.script)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("expected error %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mockResponse := `{"acknowledged": true}`

	srv, cli := testutil.NewMockServer(mockResponse, "/_scripts/my-script")
	defer srv.Close()

	shared.SetClient(cli)

	err := Delete(context.Background(), "my-script")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	srv, cli := testutil.NewMockServerWithStatus(404, `{"error":"not found"}`, "/_scripts/nonexistent")
	defer srv.Close()

	shared.SetClient(cli)

	err := Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent script")
	}

	if err.Error() != "script 'nonexistent' not found" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDeleteEmptyID(t *testing.T) {
	err := Delete(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}

	if err.Error() != "script ID is required" {
		t.Errorf("unexpected error message: %v", err)
	}
}
