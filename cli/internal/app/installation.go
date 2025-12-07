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

func InstallPackage(packageName string) error {
	// Placeholder logic for installing a package
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

	PackageSource := PackageSource{
		Name:         packageName,
		Source:       strings.TrimSpace(parsed[0]),
		Dependencies: parsed[1:],
	}

	if PackageSource.Source == "" {
		log.Fatalf("Package Installation Failed: Package not found")
		return errors.New("package source is empty")
	}

	if PackageSource.Source == "brew" {
		fmt.Printf("Installing via Homebrew: %s\n...", PackageSource.Name)
		exec.Command("brew", "install", PackageSource.Name).Run()
	} else {
		currentUser, err := user.Current()
		if err != nil {
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		fmt.Printf("Checking dependencies for package: %s\n", packageName)

		if len(PackageSource.Dependencies) > 0 {
			for _, dependency := range PackageSource.Dependencies {
				InstallPackage(dependency)
			}
		}
		fmt.Printf("Installing via Shell Command: %s\n...", PackageSource.Source)

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
		installPath := filepath.Join(cacheDir, "install_"+packageName+".sh")

		_ = os.Remove(installPath)
		os.Create(installPath)

		fileInstruct, err := os.OpenFile(installPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

		if err != nil {
			// Handle error
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		fileInstruct.Write([]byte(shellCommand + "\n"))

		defer fileInstruct.Close()

		exec.Command("chmod", "+x", installPath).Run()

		execCmd := exec.Command(installPath)

		cmdErr := execCmd.Run()
		if cmdErr != nil {
			log.Fatalf("Package Installation Failed: %v", cmdErr)
			return cmdErr
		}

		aetheisDir := filepath.Join(currentUser.HomeDir, ".aetheis")
		installedPkgsPath := filepath.Join(aetheisDir, "install_packages.json")
		file, err := os.OpenFile(installedPkgsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			// Handle error
			log.Fatalf("Package Installation Failed: %v", err)
			return err
		}

		file.Write([]byte(packageName + "\n"))

		defer file.Close()
	}

	fmt.Printf("Package %s installed successfully.\n", packageName)
	return nil
}
