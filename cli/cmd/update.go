// cmd/uninstall.go
package cmd

import (
	"fmt"
	"os"
	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [package]",
	Short: "Auto update aetheis and its packages.",
	Long: `
The 'update' command is used to update aetheis and its packages.
For example, to update a package:
  aetheis update package-name
  aetheis update
`,

	// The function that runs when the 'update' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		
		packageList := args

		fmt.Printf("Attempting to update packages: %v\n", packageList)

		err := app.Update(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update package(s): %v\n", err)
		} else {
			fmt.Printf("Successfully updated package(s)\n")
		}
	},
}