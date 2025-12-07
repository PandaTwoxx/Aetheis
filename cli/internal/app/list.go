package app

import (
	"errors"
	"fmt"
	"os/exec"
)

func ListPackages() error {
	fmt.Printf("Fetching list of installed packages...\n")
	brewPackages, err := exec.Command("brew", "list").Output()
	if err != nil {
		return errors.New("failed to fetch installed packages")
	}

	fmt.Printf("Homebrew Installed Packages:\n%s\n", string(brewPackages))

	fmt.Printf("Custom Installed Packages:\n")

	return nil
}
