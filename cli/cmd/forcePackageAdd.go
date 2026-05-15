package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var forceCmd = &cobra.Command{
	Use:   "forcePackageAdd [package]",
	Short: "Forcefully add a package to the internal Aetheis package tracker.",
	Long: `The 'forcePackageAdd' command is used to forcefully add a package to the internal Aetheis package tracker without attempting to install it.
For example, to force add a package:
  aetheis forcePackageAdd mypackage`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		err := app.ForcePackageAdd(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to force add package:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully force added package:", args[0])
	},
}
