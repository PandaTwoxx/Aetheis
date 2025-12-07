// cmd/uninstall.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [package]",
	Short: "Uninstall a package from your system.",
	Long: `
The 'uninstall' command is used to remove packages from your system.
For example, to uninstall a package:
  atheis uninstall package-name
`,

	Args: cobra.MinimumNArgs(1),

	// The function that runs when the 'uninstall' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		// Example logic for the uninstall command
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: Must specify what to uninstall (e.g., package).")
			cmd.Help()
			os.Exit(1)
		}
		
		packageList := args
		
			
		fmt.Printf("Attempting to uninstall packages: %v\n", packageList)
	},
}