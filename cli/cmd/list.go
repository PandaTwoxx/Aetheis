// cmd/install.go
package cmd

import (
	"fmt"
	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed packages on your system.",
	Long: `
The 'list' command is used to display all packages currently installed on your system.
For example, to list installed packages:
  aetheis list
`,
	Args: cobra.NoArgs,

	// The function that runs when the 'install' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		
		
			
		fmt.Printf("Listing packages\n")

		app.ListPackages()
	},
}