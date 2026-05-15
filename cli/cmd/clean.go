// cmd/clean.go
package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clear the Aetheis cache directory and remove unused dependencies.",
	Long: `
The 'clean' command performs two operations:

1. Removes all cached files stored in ~/.aetheis/cache to free up disk space
   or resolve cache-related issues.

2. Identifies and uninstalls unused dependencies - packages that were installed
   as dependencies of other packages but are no longer needed after those packages
   have been uninstalled.

This helps maintain a clean and efficient package environment.
`,

	Run: func(cmd *cobra.Command, args []string) {
		err := app.CleanCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning cache: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Cache cleaned successfully!")
	},
}
