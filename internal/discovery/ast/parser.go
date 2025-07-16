package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type ImportInfo struct {
	Position token.Position
	Alias    string
	Found    bool
}

func FindImportSpecInFile(filename, importPath string) (*ImportInfo, error) {
	fileSet := token.NewFileSet()

	file, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if isGeneratedFile(file) {
		return &ImportInfo{Found: false}, nil
	}

	for _, imp := range file.Imports {
		impPath := strings.Trim(imp.Path.Value, `"`)
		if impPath == importPath {
			info := &ImportInfo{
				Position: fileSet.Position(imp.Pos()),
				Found:    true,
			}

			if imp.Name != nil {
				info.Alias = imp.Name.Name
			}

			return info, nil
		}
	}

	return &ImportInfo{Found: false}, nil
}

func isGeneratedFile(file *ast.File) bool {
	for _, comment := range file.Comments {
		for _, c := range comment.List {
			if strings.Contains(c.Text, "Code generated") && strings.Contains(c.Text, "DO NOT EDIT") {
				return true
			}
		}
	}
	return false
}
