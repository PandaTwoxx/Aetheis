package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/PandaTwoxx/aetheis/internal/app"
	"os"
)

var loginCmd = &cobra.Command{
	Use:   "login [user] [password]",
	Short: "Login to Aetheis.",
	Long: `The 'login' command is used to log in to Aetheis.
For example, to log in:
  atheis login username password`,
	Args: cobra.MinimumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: Must specify username and password.")
			cmd.Help()
			os.Exit(1)
		}

		username := args[0]
		password := args[1]

		err := app.Login(username, password)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to log in:", err)
			os.Exit(1)
		}
		fmt.Println("Successfully logged in.")
	},
}