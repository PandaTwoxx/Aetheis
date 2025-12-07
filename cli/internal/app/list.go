package app

import (
	"fmt"
	"os/exec"
)

func ListPackages() error {
	// 1. Show Homebrew (as checked by user originally) based on system intent
	fmt.Printf("Fetching list of installed packages...\n")

	// Print managed packages from Store
	store, err := LoadPackageStore()
	if err != nil {
		fmt.Printf("Error loading package store: %v\n", err)
	} else {
		fmt.Printf("Managed Packages (Aetheis):\n")
		if len(store.Packages) == 0 {
			fmt.Printf("  (No packages installed)\n")
		}
		for _, pkg := range store.Packages {
			typeStr := "Explicit"
			if !pkg.Explicit {
				typeStr = "Dependency"
			}
			fmt.Printf("  - %s (%s)\n", pkg.Name, typeStr)
			if len(pkg.Dependencies) > 0 {
				fmt.Printf("    Dependencies: %v\n", pkg.Dependencies)
			}
		}
	}

	fmt.Println()
	// Optional: Still show brew list?
	// The original code showed "Homebrew Installed Packages" then "Custom".
	// Let's keep showing brew list as it gives context.
	brewPackages, err := exec.Command("brew", "list").Output()
	if err == nil {
		fmt.Printf("Homebrew Installed Packages (System):\n%s\n", string(brewPackages))
	}

	return nil
}
