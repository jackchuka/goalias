package discovery

import (
	"testing"

	"github.com/jackchuka/goalias/internal/discovery/ast"
)

func TestInferDefaultAlias(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		expected   string
	}{
		{
			name:       "simple package name",
			importPath: "fmt",
			expected:   "fmt",
		},
		{
			name:       "package with path",
			importPath: "github.com/user/repo",
			expected:   "repo",
		},
		{
			name:       "package with nested path",
			importPath: "github.com/user/repo/subpackage",
			expected:   "subpackage",
		},
		{
			name:       "standard library package",
			importPath: "encoding/json",
			expected:   "json",
		},
		{
			name:       "empty import path",
			importPath: "",
			expected:   "",
		},
		{
			name:       "single slash",
			importPath: "/",
			expected:   "",
		},
		{
			name:       "trailing slash",
			importPath: "github.com/user/repo/",
			expected:   "",
		},
		{
			name:       "complex path with version",
			importPath: "gopkg.in/yaml.v2",
			expected:   "yaml.v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InferDefaultAlias(tt.importPath)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetPatterns(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "empty args should return default",
			args:     []string{},
			expected: []string{"./..."},
		},
		{
			name:     "single arg",
			args:     []string{"./cmd/..."},
			expected: []string{"./cmd/..."},
		},
		{
			name:     "multiple args",
			args:     []string{"./cmd/...", "./internal/..."},
			expected: []string{"./cmd/...", "./internal/..."},
		},
		{
			name:     "nil args should return default",
			args:     nil,
			expected: []string{"./..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPatterns(tt.args)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d patterns, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("expected pattern at index %d to be %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestFindImportsInFiles(t *testing.T) {
	tests := []struct {
		name       string
		patterns   []string
		importPath string
		expectErr  bool
	}{
		{
			name:       "valid patterns and import path",
			patterns:   []string{"./..."},
			importPath: "fmt",
			expectErr:  false,
		},
		{
			name:       "empty patterns with import path",
			patterns:   []string{},
			importPath: "fmt",
			expectErr:  false,
		},
		{
			name:       "nonexistent import path",
			patterns:   []string{"./..."},
			importPath: "nonexistent/package",
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindImportsInFiles(tt.patterns, tt.importPath)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("expected non-nil result")
				return
			}

			// For nonexistent import paths, we expect an empty slice
			if tt.importPath == "nonexistent/package" {
				if len(result) != 0 {
					t.Errorf("expected empty result for nonexistent import path, got %d results", len(result))
				}
				return
			}

			// Validate the structure of results
			for i, importResult := range result {
				if importResult.File == "" {
					t.Errorf("result[%d].File should not be empty", i)
				}
				if importResult.Location == "" {
					t.Errorf("result[%d].Location should not be empty", i)
				}
				if importResult.Alias == "" {
					t.Errorf("result[%d].Alias should not be empty", i)
				}
				if importResult.Info == nil {
					t.Errorf("result[%d].Info should not be nil", i)
				}
				if !importResult.Info.Found {
					t.Errorf("result[%d].Info.Found should be true", i)
				}
			}
		})
	}
}

func TestImportResult(t *testing.T) {
	// Test the ImportResult struct creation and field access
	info := &ast.ImportInfo{
		Found: true,
		Alias: "f",
	}

	result := ImportResult{
		File:     "/path/to/file.go",
		Location: "/path/to/file.go:10",
		Alias:    "f",
		Info:     info,
	}

	if result.File != "/path/to/file.go" {
		t.Errorf("expected File to be '/path/to/file.go', got %q", result.File)
	}

	if result.Location != "/path/to/file.go:10" {
		t.Errorf("expected Location to be '/path/to/file.go:10', got %q", result.Location)
	}

	if result.Alias != "f" {
		t.Errorf("expected Alias to be 'f', got %q", result.Alias)
	}

	if result.Info != info {
		t.Errorf("expected Info to be the same instance")
	}

	if !result.Info.Found {
		t.Errorf("expected Info.Found to be true")
	}

	if result.Info.Alias != "f" {
		t.Errorf("expected Info.Alias to be 'f', got %q", result.Info.Alias)
	}
}
