package app

import (
	"os/exec"
	"os/user"
	"os"
	"errors"
	"fmt"
	"log"
)

func InitializeEnvironment() error {
	fmt.Printf("Setting up Aetheis environment...\n")

	// Example: Check if Homebrew is installed
	_, err := exec.LookPath("brew")
	if err != nil {
		log.Fatalf("Homebrew is not installed. Please install Homebrew with aetheis install brew or aetheis install brew-local.")
		return errors.New("homebrew not installed")
	}

	// Create directories + files

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
		return err
	}

	aetheisDir := fmt.Sprintf("%s/.aetheis", currentUser.HomeDir)
	cacheDir := fmt.Sprintf("%s/cache", aetheisDir)

	err = exec.Command("mkdir", "-p", cacheDir).Run()
	if err != nil {
		log.Fatalf("Failed to create Aetheis directories: %v", err)
		return err
	}

	os.Create(fmt.Sprintf("%s/install_packages.json", aetheisDir))

	fmt.Printf("Aetheis environment initialized successfully at %s\n", aetheisDir)
	return nil
}