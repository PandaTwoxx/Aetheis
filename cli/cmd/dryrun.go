// cmd/dryrun.go
package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var dryrunCmd = &cobra.Command{
	Use:   "dryrun [action] [package]",
	Short: "Preview commands that will run for install/uninstall/update/upgrade operations.",
	Long: `
The 'dryrun' command shows you exactly what commands will be executed when you install, uninstall, update, or upgrade packages.
This is useful for understanding what changes will be made to your system before actually running them.

Usage:
  aetheis dryrun install package-name
  aetheis dryrun uninstall package-name
  aetheis dryrun update
  aetheis dryrun upgrade [package-name]

Examples:
  aetheis dryrun install nodejs
  aetheis dryrun uninstall python
  aetheis dryrun update
  aetheis dryrun upgrade
  aetheis dryrun upgrade ruby
`,
	Args: cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Error: Must specify an action (install, uninstall, update, or upgrade).")
			cmd.Help()
			os.Exit(1)
		}

		action := args[0]
		packageList := args[1:]

		switch action {
		case "install":
			if len(packageList) == 0 {
				fmt.Fprintln(os.Stderr, "Error: Must specify a package to install.")
				cmd.Help()
				os.Exit(1)
			}
			fmt.Printf("=== DRY RUN: Install Commands for %v ===\n\n", packageList)
			for _, pkg := range packageList {
				err := app.PreviewInstallPackage(pkg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error previewing install for %s: %v\n", pkg, err)
				}
			}

		case "uninstall":
			if len(packageList) == 0 {
				fmt.Fprintln(os.Stderr, "Error: Must specify a package to uninstall.")
				cmd.Help()
				os.Exit(1)
			}
			fmt.Printf("=== DRY RUN: Uninstall Commands for %v ===\n\n", packageList)
			for _, pkg := range packageList {
				err := app.PreviewUninstallPackage(pkg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error previewing uninstall for %s: %v\n", pkg, err)
				}
			}

		case "update":
			fmt.Printf("=== DRY RUN: Update aetheis ===\n")
			err := app.PreviewUpdate(packageList)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error previewing update: %v\n", err)
			}

		case "upgrade":
			fmt.Printf("=== DRY RUN: Upgrade Commands for %v ===\n", packageList)
			err := app.PreviewUpgrade(packageList)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error previewing upgrade: %v\n", err)
			}

		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown action '%s'. Must be 'install', 'uninstall', 'update', or 'upgrade'.\n", action)
			cmd.Help()
			os.Exit(1)
		}
	},
}
