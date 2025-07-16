package commands

import (
	"fmt"
	"os"

	"github.com/jackchuka/goalias/internal/discovery"
	"github.com/jackchuka/goalias/internal/lsp"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [packages]",
	Short: "Set import alias for a package",
	Long: `Set import alias for a package across specified Go packages.
	
Examples:
  goalias set -p github.com/example/mypackage -a mypkg
  goalias set -p github.com/example/mypackage -a mypkg ./cmd/...`,
	RunE: runSet,
}

var (
	setPackage string
	setAlias   string
	setPreview bool
)

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&setPackage, "package", "p", "", "Full import path to manage (required)")
	setCmd.Flags().StringVarP(&setAlias, "alias", "a", "", "Desired alias identifier (required)")
	setCmd.Flags().BoolVarP(&setPreview, "preview", "n", false, "Show diff instead of writing changes")

	_ = setCmd.MarkFlagRequired("package")
	_ = setCmd.MarkFlagRequired("alias")
}

func runSet(cmd *cobra.Command, args []string) error {
	patterns := discovery.GetPatterns(args)

	results, err := discovery.FindImportsInFiles(patterns, setPackage)
	if err != nil {
		return err
	}

	var filesToProcess []discovery.ImportResult

	for _, result := range results {
		// Use the effective alias (which includes inferred default aliases)
		if result.Alias == setAlias {
			continue
		}
		filesToProcess = append(filesToProcess, result)
	}

	if len(filesToProcess) == 0 {
		fmt.Println("No files need updating")
		return nil
	}

	// Get current working directory for LSP client
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Create LSP client
	client, err := lsp.NewClient(cwd)
	if err != nil {
		return fmt.Errorf("failed to create LSP client: %w", err)
	}
	defer func() {
		_ = client.Close()
	}()

	// Initialize LSP client
	if err := client.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize LSP client: %w", err)
	}

	fmt.Printf("Processing %d files...\n", len(filesToProcess))

	// Process files using LSP client
	for i, result := range filesToProcess {
		fmt.Printf("Processing file %d/%d: %s\n", i+1, len(filesToProcess), result.File)

		if err := processFileWithLSP(client, result); err != nil {
			return fmt.Errorf("failed to process %s: %w", result.File, err)
		}
	}

	return nil
}

func processFileWithLSP(client *lsp.Client, result discovery.ImportResult) error {
	// Convert Go token position to LSP position (0-based)
	line := result.Info.Position.Line - 1
	column := result.Info.Position.Column - 1

	// Perform rename operation
	workspaceEdit, err := client.Rename(result.File, line, column, setAlias)
	if err != nil {
		return fmt.Errorf("rename operation failed: %w", err)
	}

	// Apply the workspace edit
	if err := lsp.ApplyWorkspaceEdit(workspaceEdit, setPreview); err != nil {
		return fmt.Errorf("failed to apply workspace edit: %w", err)
	}

	return nil
}
