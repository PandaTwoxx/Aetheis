package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type PackageSource struct {
	Name         string
	Source       string
	Dependencies []string
}

func InstallPackage(packageName string, explicit bool) error {
	// 1. Load Store
	store, err := LoadPackageStore()
	if err != nil {
		log.Fatalf("Failed to load package store: %v", err)
		return err
	}

	// 2. Check if already installed (optional, but good optimization)
	// For now, we reinstall to ensure latest version or if something broke,
	// but we should respect the 'explicit' flag update if it was previously implicit.
	if existing := store.GetPackage(packageName); existing != nil {
		if explicit && !existing.Explicit {
			existing.Explicit = true
			SavePackageStore(store)
		}
		fmt.Printf("Package %s is already installed.\n", packageName)
		// We could return nil here, but maybe we want to verify dependencies again?
		// Let's assume re-install is okay or we just return.
		// For matching original behavior regarding dependencies, we might want to check them.
		// But let's proceed with install to be safe.
	}

	if packageName == "" {
		return errors.New("package name cannot be empty")
	}
	fmt.Printf("Installing package: %s\n", packageName)

	resp, err := http.Get("https://aetheis.vercel.app/" + packageName)

	if err != nil {
		log.Fatalf("Package Installation Failed: %v", err)
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Package Installation Failed: %v", err)
		return err
	}

	parsed := strings.Split(string(bodyBytes), " ")

	targetPackage := InstalledPackage{
		Name:         packageName,
		Explicit:     explicit,
		Dependencies: []string{},
	}

	sourceCode := strings.TrimSpace(parsed[0])
	if len(parsed) > 1 {
		targetPackage.Dependencies = parsed[1:]
	}

	if sourceCode == "" {
		log.Fatalf("Package Installation Failed: Package not found")
		return errors.New("package source is empty")
	}

	if sourceCode == "brew" {
		fmt.Printf("Installing via Homebrew: %s\n...", targetPackage.Name)
		exec.Command("brew", "install", targetPackage.Name).Run()
	} else {
		currentUser, err := user.Current()
		if err != nil {
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		fmt.Printf("Checking dependencies for package: %s\n", packageName)

		if len(targetPackage.Dependencies) > 0 {
			for _, dependency := range targetPackage.Dependencies {
				// Recursive install with explicit=false
				InstallPackage(dependency, false)
			}
		}
		fmt.Printf("Installing via Shell Command: %s\n...", sourceCode)

		resp, err := http.Get("https://aetheis.vercel.app/install/" + packageName)

		if err != nil {
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		shellCommand := string(bodyBytes)

		if shellCommand == "" {
			log.Fatalf("Package Installation Failed: Install command not found")
			return errors.New("install command is empty")
		}

		cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
		_ = os.MkdirAll(cacheDir, 0755) // Ensure cache dir exists
		installPath := filepath.Join(cacheDir, "install_"+packageName+".sh")

		_ = os.Remove(installPath)
		// No need to Create then Open, just WriteFile or OpenFile with Create
		fileInstruct, err := os.OpenFile(installPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

		if err != nil {
			// Handle error
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		// Prepend shebang as per previous fix
		fileInstruct.Write([]byte("#!/bin/sh\n" + shellCommand + "\n"))

		fileInstruct.Close()

		exec.Command("chmod", "+x", installPath).Run()

		execCmd := exec.Command("/bin/sh", installPath)
		// Connect stdout/stderr to see output
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		cmdErr := execCmd.Run()
		if cmdErr != nil {
			log.Fatalf("Package Installation Failed: %v", cmdErr)
			return cmdErr
		}
	}

	// Save to Store
	store.AddPackage(targetPackage)
	if err := SavePackageStore(store); err != nil {
		log.Printf("Warning: Failed to save package store: %v", err)
	}

	fmt.Printf("Package %s installed successfully.\n", packageName)
	return nil
}
