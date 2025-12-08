package cmd

import (
	"github.com/spf13/cobra"
	"github.com/PandaTwoxx/aetheis/internal/app"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Signup to Aetheis",
	Run: func(cmd *cobra.Command, args []string) {
		app.Signup(args[0], args[1])
	},
}
