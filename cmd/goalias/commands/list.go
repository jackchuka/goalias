package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jackchuka/goalias/internal/discovery"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [packages]",
	Short: "List import aliases for a package",
	Long: `List all files and their aliases for a package within specified Go packages.
	
Examples:
  goalias list -p github.com/example/mypackage
  goalias list -p github.com/example/mypackage ./cmd/...`,
	RunE: runList,
}

var listPackage string

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listPackage, "package", "p", "", "Full import path to manage (required)")

	_ = listCmd.MarkFlagRequired("package")
}

func runList(cmd *cobra.Command, args []string) error {
	patterns := discovery.GetPatterns(args)

	results, err := discovery.FindImportsInFiles(patterns, listPackage)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Printf("No imports found for package: %s\n", listPackage)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "LOCATION\tALIAS")
	_, _ = fmt.Fprintln(w, "--------\t-----")

	for _, r := range results {
		_, _ = fmt.Fprintf(w, "%s\t%s\n", r.Location, r.Alias)
	}

	return w.Flush()
}
