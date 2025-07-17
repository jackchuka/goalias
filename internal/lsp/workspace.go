package lsp

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"go.lsp.dev/protocol"
)

// ApplyWorkspaceEdit applies a workspace edit to the filesystem
func ApplyWorkspaceEdit(edit *protocol.WorkspaceEdit, preview bool) error {
	if edit == nil {
		return fmt.Errorf("workspace edit is nil")
	}

	// Handle changes map (deprecated but still used)
	if edit.Changes != nil {
		for uri, edits := range edit.Changes {
			if err := applyTextEdits(string(uri), edits, preview); err != nil {
				return fmt.Errorf("failed to apply changes to %s: %w", uri, err)
			}
		}
	}

	// Handle document changes (preferred method)
	if edit.DocumentChanges != nil {
		for _, docChange := range edit.DocumentChanges {
			// docChange is already a protocol.TextDocumentEdit, no need to type assert
			if err := applyTextEdits(string(docChange.TextDocument.URI), docChange.Edits, preview); err != nil {
				return fmt.Errorf("failed to apply document changes to %s: %w", docChange.TextDocument.URI, err)
			}
		}
	}

	return nil
}

// applyTextEdits applies text edits to a file
func applyTextEdits(uri string, edits []protocol.TextEdit, preview bool) error {
	// Convert URI to file path
	filePath := uriToFilePath(uri)

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Apply edits
	modifiedContent, err := applyEditsToContent(string(content), edits)
	if err != nil {
		return fmt.Errorf("failed to apply edits: %w", err)
	}

	if preview {
		// Print diff-like output
		fmt.Printf("--- %s\n", filePath)
		fmt.Printf("+++ %s\n", filePath)

		// Simple diff output (could be enhanced with proper diff library)
		originalLines := strings.Split(string(content), "\n")
		modifiedLines := strings.Split(modifiedContent, "\n")

		for i, line := range originalLines {
			if i < len(modifiedLines) {
				if line != modifiedLines[i] {
					fmt.Printf("-%s\n", line)
					fmt.Printf("+%s\n", modifiedLines[i])
				}
			} else {
				fmt.Printf("-%s\n", line)
			}
		}

		// Print any additional lines in modified content
		for i := len(originalLines); i < len(modifiedLines); i++ {
			fmt.Printf("+%s\n", modifiedLines[i])
		}

		fmt.Println()
	} else {
		// Write the modified content back to the file
		if err := os.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

// applyEditsToContent applies text edits to content string
func applyEditsToContent(content string, edits []protocol.TextEdit) (string, error) {
	if len(edits) == 0 {
		return content, nil
	}

	// Split content into lines for easier manipulation
	lines := strings.Split(content, "\n")

	// Sort edits by position (reverse order: end to start) to avoid offset issues
	// when multiple edits occur on the same line
	sortedEdits := make([]protocol.TextEdit, len(edits))
	copy(sortedEdits, edits)

	// Sort in reverse order: later positions first, then earlier positions
	// This ensures that character positions remain valid as we apply edits
	sort.Slice(sortedEdits, func(i, j int) bool {
		edit1 := sortedEdits[i]
		edit2 := sortedEdits[j]

		// Sort by line first (later lines first)
		if edit1.Range.Start.Line != edit2.Range.Start.Line {
			return edit1.Range.Start.Line > edit2.Range.Start.Line
		}

		// For same line, sort by character position (later positions first)
		return edit1.Range.Start.Character > edit2.Range.Start.Character
	})

	// Apply edits in reverse position order
	for _, edit := range sortedEdits {
		lines = applyEditToLines(lines, edit)
	}

	return strings.Join(lines, "\n"), nil
}

// applyEditToLines applies a single text edit to lines
func applyEditToLines(lines []string, edit protocol.TextEdit) []string {
	startLine := int(edit.Range.Start.Line)
	startChar := int(edit.Range.Start.Character)
	endLine := int(edit.Range.End.Line)
	endChar := int(edit.Range.End.Character)

	// Handle single line edit
	if startLine == endLine {
		if startLine < len(lines) {
			line := lines[startLine]
			if startChar <= len(line) && endChar <= len(line) {
				newLine := line[:startChar] + edit.NewText + line[endChar:]
				lines[startLine] = newLine
			}
		}
		return lines
	}

	// Handle multi-line edit
	if startLine < len(lines) && endLine < len(lines) {
		// Get the prefix from start line
		prefix := ""
		if startChar < len(lines[startLine]) {
			prefix = lines[startLine][:startChar]
		}

		// Get the suffix from end line
		suffix := ""
		if endChar < len(lines[endLine]) {
			suffix = lines[endLine][endChar:]
		}

		// Create new content
		newContent := prefix + edit.NewText + suffix

		// Replace the range with new content
		result := make([]string, 0, len(lines)-(endLine-startLine))
		result = append(result, lines[:startLine]...)
		result = append(result, newContent)
		result = append(result, lines[endLine+1:]...)

		return result
	}

	return lines
}

// uriToFilePath converts a file URI to a file path
func uriToFilePath(uri string) string {
	// Remove file:// prefix
	if strings.HasPrefix(uri, "file://") {
		return uri[7:]
	}
	return uri
}
