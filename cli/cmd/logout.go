package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/PandaTwoxx/aetheis/internal/app"
	"os"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of Aetheis.",
	Long: `The 'logout' command is used to log out of Aetheis.
For example, to log out:
  atheis logout`,

	Run: func(cmd *cobra.Command, args []string) {
		err := app.Logout()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to log out:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully logged out.")
	},
}