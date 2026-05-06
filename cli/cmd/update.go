// cmd/uninstall.go
package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update aetheis itself.",
	Long: `
The 'update' command updates aetheis to the latest version using the registry.
This updates the aetheis package manager itself, not the packages you have installed.

To upgrade your installed packages, use the 'upgrade' command instead:
  aetheis upgrade
  aetheis upgrade package-name
`,

	// The function that runs when the 'update' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔄 Starting aetheis update...")

		err := app.Update(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to update aetheis: %v\n", err)
			os.Exit(1)
		}
	},
}
