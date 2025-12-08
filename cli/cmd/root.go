// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Define a variable to hold the root command
var rootCmd = &cobra.Command{
	Use:   "aetheis",
	Short: "An absolutely marvelous system package manager.",
	Long: `
This is a package manager for your system. Its simple yet complex and powerful.
You can use it to install, remove, and manage packages on your system from multiple sources including homebrew, local files, and our own custom repositories.
Usage is simple and intuitive, making it easy for both beginners and advanced users to manage their software.
`,
	// This function runs if the command is called without any subcommands.
	Run: func(cmd *cobra.Command, args []string) {
		// By default, just print the help message
		cmd.Help()
	},
}

// Execute is the main entry point called from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	// 2. Define global flags that apply to ALL commands
	rootCmd.PersistentFlags().StringP("config", "c", "", "Config package installation path (default is $HOME/.aetheis/)")

	// 3. Define local flags that only apply to the root command (optional)
	// rootCmd.Flags().BoolP("version", "v", false, "Print version number")
}
