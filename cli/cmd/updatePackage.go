package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/PandaTwoxx/aetheis/internal/app"
)

var updatePackageCmd = &cobra.Command{
	Use:   "updatePackage",
	Short: "Update a package in Aetheis.",
	Long: `The 'updatePackage' command is used to update a package in Aetheis.
For example, to update a package:
  atheis updatePackage <package>`,

	Run: func(cmd *cobra.Command, args []string) {
		err := app.UpdatePackage(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to update package:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully updated package.")
	},
}