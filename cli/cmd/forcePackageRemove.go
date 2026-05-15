package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var forcePackageRemoveCmd = &cobra.Command{
	Use:   "forcePackageRemove [package]",
	Short: "Forcefully remove a package from Aetheis without dependency checks.",
	Long: `The 'forcePackageRemove' command is used to forcefully remove a package from Aetheis
without checking if other packages depend on it. This bypasses the normal dependency
validation and removes the package immediately from the system.

WARNING: Use with caution! Removing a package that is required by other packages
may break those packages. Only use this command if you know what you're doing.

For example, to force remove a package:
  aetheis forcePackageRemove mypackage`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		err := app.ForcePackageRemove(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to force remove package:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully force removed package:", args[0])
	},
}
