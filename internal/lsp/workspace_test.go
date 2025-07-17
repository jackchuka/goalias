package lsp

import (
	"testing"

	"go.lsp.dev/protocol"
)

func TestUriToFilePath(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "file URI with protocol",
			uri:      "file:///Users/test/file.go",
			expected: "/Users/test/file.go",
		},
		{
			name:     "file URI without protocol",
			uri:      "/Users/test/file.go",
			expected: "/Users/test/file.go",
		},
		{
			name:     "Windows file URI",
			uri:      "file:///C:/Users/test/file.go",
			expected: "/C:/Users/test/file.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uriToFilePath(tt.uri)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestApplyEditToLines(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		edit     protocol.TextEdit
		expected []string
	}{
		{
			name:  "single line replacement",
			lines: []string{"import \"fmt\"", "func main() {}"},
			edit: protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 8},
					End:   protocol.Position{Line: 0, Character: 11},
				},
				NewText: "os",
			},
			expected: []string{"import \"os\"", "func main() {}"},
		},
		{
			name:  "single line insertion",
			lines: []string{"import \"fmt\"", "func main() {}"},
			edit: protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 7},
					End:   protocol.Position{Line: 0, Character: 7},
				},
				NewText: "alias ",
			},
			expected: []string{"import alias \"fmt\"", "func main() {}"},
		},
		{
			name:  "multi-line replacement",
			lines: []string{"import (", "\"fmt\"", "\"os\"", ")", "func main() {}"},
			edit: protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: 1, Character: 0},
					End:   protocol.Position{Line: 2, Character: 4},
				},
				NewText: "f \"fmt\"",
			},
			expected: []string{"import (", "f \"fmt\"", ")", "func main() {}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyEditToLines(tt.lines, tt.edit)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d lines, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("line %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestApplyEditsToContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		edits    []protocol.TextEdit
		expected string
	}{
		{
			name:     "no edits",
			content:  "package main\n\nimport \"fmt\"\n\nfunc main() {}",
			edits:    []protocol.TextEdit{},
			expected: "package main\n\nimport \"fmt\"\n\nfunc main() {}",
		},
		{
			name:    "single edit",
			content: "package main\n\nimport \"fmt\"\n\nfunc main() {}",
			edits: []protocol.TextEdit{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 2, Character: 8},
						End:   protocol.Position{Line: 2, Character: 11},
					},
					NewText: "os",
				},
			},
			expected: "package main\n\nimport \"os\"\n\nfunc main() {}",
		},
		{
			name:    "multiple edits on same line - config to pkg_config",
			content: "environment == config.EnvironmentDev || environment == config.EnvironmentProd,",
			edits: []protocol.TextEdit{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 15},
						End:   protocol.Position{Line: 0, Character: 21},
					},
					NewText: "pkg_config",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 55},
						End:   protocol.Position{Line: 0, Character: 61},
					},
					NewText: "pkg_config",
				},
			},
			expected: "environment == pkg_config.EnvironmentDev || environment == pkg_config.EnvironmentProd,",
		},
		{
			name:    "multiple edits - import statements",
			content: "import \"fmt\"; import \"os\"; import \"log\"",
			edits: []protocol.TextEdit{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 8},
						End:   protocol.Position{Line: 0, Character: 11},
					},
					NewText: "myfmt",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 22},
						End:   protocol.Position{Line: 0, Character: 24},
					},
					NewText: "myos",
				},
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 35},
						End:   protocol.Position{Line: 0, Character: 38},
					},
					NewText: "mylog",
				},
			},
			expected: "import \"myfmt\"; import \"myos\"; import \"mylog\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyEditsToContent(tt.content, tt.edits)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
