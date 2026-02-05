package cmd

import (
	"fmt"
	"os"

	"github.com/PandaTwoxx/aetheis/internal/app"
	"github.com/spf13/cobra"
)

var listRegistryCmd = &cobra.Command{
	Use:   "list-registry",
	Short: "List all packages on the Aetheis registry.",
	Long: `List all packages available on the Aetheis package registry.
Supports filtering by name or owner (partial match, case-insensitive).
For example:
  aetheis list-registry
  aetheis list-registry --name deno
  aetheis list-registry --owner panda`,
	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		owner, _ := cmd.Flags().GetString("owner")
		if err := app.ListRegistryPackages(name, owner); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	listRegistryCmd.Flags().StringP("name", "n", "", "Filter packages by name (partial match)")
	listRegistryCmd.Flags().StringP("owner", "o", "", "Filter packages by owner (partial match)")
}

