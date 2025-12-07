// cmd/install.go
package cmd

import (
	"fmt"
	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Aetheis environment.",
	Long: `
The 'init' command is used to set up the initial environment for Aetheis.
For example, to initialize the environment:
  aetheis init
`,
	Args: cobra.NoArgs,

	// The function that runs when the 'install' command is executed.
	Run: func(cmd *cobra.Command, args []string) {
		
		
			
		fmt.Printf("Initializing Aetheis environment\n")

		app.InitializeEnvironment()
	},
}