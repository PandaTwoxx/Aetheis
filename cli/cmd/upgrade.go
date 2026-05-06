// cmd/upgrade.go
package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [package]",
	Short: "Upgrade installed packages.",
	Long: `
The 'upgrade' command updates your installed packages to their latest versions.
You can upgrade all packages or specify individual packages to upgrade.

Examples:
  aetheis upgrade              (upgrade all packages)
  aetheis upgrade nodejs       (upgrade nodejs)
  aetheis upgrade python ruby  (upgrade multiple packages)
`,

	// The function that runs when the 'upgrade' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("🔄 Upgrading all installed packages...")
		} else {
			fmt.Printf("🔄 Upgrading packages: %v\n", args)
		}

		err := app.Upgrade(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to upgrade package(s): %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("✅ Successfully upgraded package(s)\n")
		}
	},
}
