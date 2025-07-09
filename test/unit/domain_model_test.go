package unit

import (
	"encoding/json"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/domain/model"
)

func TestScaffoldOptions(t *testing.T) {
	// Test JSON marshaling and unmarshaling
	opts := model.ScaffoldOptions{
		AppType:      "api",
		DatabaseType: "postgresql",
		RouterType:   "echo",
		ConfigType:   "env",
		LogFormat:    "json",
		ModulePath:   "github.com/testuser/testproject",
		Features:     []string{"migrations", "logging"},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("Failed to marshal ScaffoldOptions: %v", err)
	}

	// Unmarshal back
	var unmarshalledOpts model.ScaffoldOptions
	err = json.Unmarshal(jsonData, &unmarshalledOpts)
	if err != nil {
		t.Fatalf("Failed to unmarshal ScaffoldOptions: %v", err)
	}

	// Compare fields
	if opts.AppType != unmarshalledOpts.AppType {
		t.Errorf("AppType does not match: expected %s, got %s", opts.AppType, unmarshalledOpts.AppType)
	}
	if opts.DatabaseType != unmarshalledOpts.DatabaseType {
		t.Errorf("DatabaseType does not match: expected %s, got %s", opts.DatabaseType, unmarshalledOpts.DatabaseType)
	}
	if opts.RouterType != unmarshalledOpts.RouterType {
		t.Errorf("RouterType does not match: expected %s, got %s", opts.RouterType, unmarshalledOpts.RouterType)
	}
	if opts.ModulePath != unmarshalledOpts.ModulePath {
		t.Errorf("ModulePath does not match: expected %s, got %s", opts.ModulePath, unmarshalledOpts.ModulePath)
	}

	// Check features array
	if len(opts.Features) != len(unmarshalledOpts.Features) {
		t.Errorf("Features length does not match: expected %d, got %d",
			len(opts.Features), len(unmarshalledOpts.Features))
	} else {
		for i, feature := range opts.Features {
			if feature != unmarshalledOpts.Features[i] {
				t.Errorf("Feature at index %d does not match: expected %s, got %s",
					i, feature, unmarshalledOpts.Features[i])
			}
		}
	}
}

func TestGeneratedScaffold(t *testing.T) {
	// Test struct initialization
	scaffold := model.GeneratedScaffold{
		ID: "test-id",
		Options: model.ScaffoldOptions{
			AppType:      "api",
			DatabaseType: "mysql",
			RouterType:   "chi",
		},
		CreatedAt: "2023-06-01T12:00:00Z",
		FilePath:  "/tmp/test-scaffold.zip",
		Size:      12345,
	}

	// Verify fields
	if scaffold.ID != "test-id" {
		t.Errorf("ID does not match: expected test-id, got %s", scaffold.ID)
	}
	if scaffold.Options.AppType != "api" {
		t.Errorf("Options.AppType does not match: expected api, got %s", scaffold.Options.AppType)
	}
	if scaffold.Options.DatabaseType != "mysql" {
		t.Errorf("Options.DatabaseType does not match: expected mysql, got %s", scaffold.Options.DatabaseType)
	}
	if scaffold.FilePath != "/tmp/test-scaffold.zip" {
		t.Errorf("FilePath does not match: expected /tmp/test-scaffold.zip, got %s", scaffold.FilePath)
	}
	if scaffold.Size != 12345 {
		t.Errorf("Size does not match: expected 12345, got %d", scaffold.Size)
	}
}
