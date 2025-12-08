package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/PandaTwoxx/aetheis/internal/app"
)

var deletePackageCmd = &cobra.Command{
	Use:   "deletePackage",
	Short: "Delete a package from Aetheis.",
	Long: `The 'deletePackage' command is used to delete a package from Aetheis.
For example, to delete a package:
  atheis deletePackage <package>`,

	Run: func(cmd *cobra.Command, args []string) {
		err := app.DeletePackage(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to delete package:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully deleted package.")
	},
}
