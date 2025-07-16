package discovery

import (
	"fmt"
	"strings"

	"github.com/jackchuka/goalias/internal/discovery/ast"
)

type ImportResult struct {
	File     string
	Location string
	Alias    string
	Info     *ast.ImportInfo
}

func FindImportsInFiles(patterns []string, importPath string) ([]ImportResult, error) {
	packages, err := ListPackages(patterns)
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}

	files := GetGoFilesFromPackages(packages)
	results := make([]ImportResult, 0)

	for _, file := range files {
		info, err := ast.FindImportSpecInFile(file, importPath)
		if err != nil {
			continue
		}

		if !info.Found {
			continue
		}

		alias := info.Alias
		if alias == "" {
			alias = InferDefaultAlias(importPath)
		}

		location := fmt.Sprintf("%s:%d", file, info.Position.Line)
		results = append(results, ImportResult{
			File:     file,
			Location: location,
			Alias:    alias,
			Info:     info,
		})
	}

	return results, nil
}

func GetPatterns(args []string) []string {
	patterns := []string{"./..."}
	if len(args) > 0 {
		patterns = args
	}
	return patterns
}

func InferDefaultAlias(importPath string) string {
	parts := strings.Split(importPath, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
