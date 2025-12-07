// cmd/install.go
package cmd

import (
	"fmt"
	"os"
	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Install a new package on your system.",
	Long: `
The 'install' command is used to add new packages to your system.
For example, to install a package:
  atheis install package-name
`,

	Args: cobra.MinimumNArgs(1),

	// The function that runs when the 'install' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		// Example logic for the install command
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: Must specify what to install (e.g., package).")
			cmd.Help()
			os.Exit(1)
		}
		
		packageList := args
		
			
		fmt.Printf("Attempting to install packages: %v\n", packageList)

		for _, pkg := range packageList {
			err := app.InstallPackage(pkg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to install package %s: %v\n", pkg, err)
			} else {
				fmt.Printf("Successfully installed package: %s\n", pkg)
			}
		}
	},
}