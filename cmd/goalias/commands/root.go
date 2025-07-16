package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goalias",
	Short: "Standardize Go import aliases across your codebase",
	Long: `Standardize Go import aliases across your codebase
goalias automatically manages import aliases in Go projects, ensuring consistency across all files using the Go Language Server (gopls) for fast, accurate refactoring.`,
}

func Execute() error {
	return rootCmd.Execute()
}
