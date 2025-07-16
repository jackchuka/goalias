package discovery

import (
	"testing"
)

func TestGetGoFilesFromPackages(t *testing.T) {
	tests := []struct {
		name     string
		packages []Package
		expected []string
	}{
		{
			name:     "empty packages",
			packages: []Package{},
			expected: []string{},
		},
		{
			name: "single package with one file",
			packages: []Package{
				{
					ImportPath: "example.com/test",
					Dir:        "/path/to/test",
					GoFiles:    []string{"main.go"},
				},
			},
			expected: []string{"/path/to/test/main.go"},
		},
		{
			name: "single package with multiple files",
			packages: []Package{
				{
					ImportPath: "example.com/test",
					Dir:        "/path/to/test",
					GoFiles:    []string{"main.go", "helper.go", "types.go"},
				},
			},
			expected: []string{
				"/path/to/test/main.go",
				"/path/to/test/helper.go",
				"/path/to/test/types.go",
			},
		},
		{
			name: "multiple packages with files",
			packages: []Package{
				{
					ImportPath: "example.com/test",
					Dir:        "/path/to/test",
					GoFiles:    []string{"main.go"},
				},
				{
					ImportPath: "example.com/utils",
					Dir:        "/path/to/utils",
					GoFiles:    []string{"helper.go", "types.go"},
				},
			},
			expected: []string{
				"/path/to/test/main.go",
				"/path/to/utils/helper.go",
				"/path/to/utils/types.go",
			},
		},
		{
			name: "package with no go files",
			packages: []Package{
				{
					ImportPath: "example.com/test",
					Dir:        "/path/to/test",
					GoFiles:    []string{},
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGoFilesFromPackages(tt.packages)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d files, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("expected file at index %d to be %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestListPackages(t *testing.T) {
	tests := []struct {
		name      string
		patterns  []string
		expectErr bool
	}{
		{
			name:      "empty patterns should use default",
			patterns:  []string{},
			expectErr: false,
		},
		{
			name:      "single pattern",
			patterns:  []string{"./..."},
			expectErr: false,
		},
		{
			name:      "multiple patterns",
			patterns:  []string{"./cmd/...", "./internal/..."},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test will actually run 'go list' command
			// In a real project, we might want to mock this or skip in CI
			result, err := ListPackages(tt.patterns)

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
			}
		})
	}
}

func TestListPackagesWithInvalidPattern(t *testing.T) {
	// Test with an invalid pattern that should cause go list to fail
	patterns := []string{"./nonexistent/..."}
	result, err := ListPackages(patterns)

	// This should not necessarily fail as go list might return empty results
	// for non-existent patterns, but the function should handle it gracefully
	if err != nil {
		// Error is acceptable for invalid patterns
		if result != nil {
			t.Errorf("expected nil result when error occurs")
		}
	}
}
