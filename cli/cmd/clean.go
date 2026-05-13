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
	Short: "Clear the Aetheis cache directory.",
	Long: `
The 'clean' command removes all cached files stored in ~/.aetheis/cache.
This can help free up disk space or resolve cache-related issues.
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
