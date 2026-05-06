package cmd

import (
	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup [username] [password]",
	Short: "Signup to Aetheis",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.Signup(args[0], args[1])
	},
}
