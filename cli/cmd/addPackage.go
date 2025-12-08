package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/PandaTwoxx/aetheis/internal/app"
)

var addPackageCmd = &cobra.Command{
	Use:   "addPackage [package]",
	Short: "Add a package to Aetheis.",
	Long: `The 'addPackage' command is used to add a package to Aetheis.
For example, to add a package:
  atheis addPackage <package>`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		err := app.AddPackage(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to add package:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully added package.")
	},
}