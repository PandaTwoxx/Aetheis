package app

import (
	"os/exec"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func UninstallPackage(packageName string) error {
	// Placeholder logic for uninstalling a package
	if packageName == "" {
		return errors.New("package name cannot be empty")
	}
	fmt.Printf("Uninstalling package: %s\n", packageName)

	resp, err := http.Get("https://aetheis.vercel.app/" + packageName)

	if err != nil {
		log.Fatalf("Package Uninstallation Failed: %v", err)
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Package Uninstallation Failed: %v", err)
	}

	PackageSource := PackageSource{
		Name:   packageName,
		Source: string(bodyBytes),
	}

	if PackageSource.Source == "" {
		log.Fatalf("Package Uninstallation Failed: Package not found")
		return errors.New("package source is empty")
	}

	if PackageSource.Source == "brew"{
		fmt.Printf("Uninstalling via Homebrew: %s\n...", PackageSource.Name)
		exec.Command("brew", "uninstall", PackageSource.Name).Run()
	} else{
		fmt.Printf("Uninstalling via Shell Command: %s\n...", PackageSource.Source)
		
		resp, err := http.Get("https://aetheis.vercel.app/uninstall/" + packageName)

		if err != nil {
			log.Fatalf("Package Uninstallation Failed: %v", err)
			return err
		}

		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Package Uninstallation Failed: %v", err)
		}

		shellCommand := string(bodyBytes)

		if shellCommand == "" {
			log.Fatalf("Package Uninstallation Failed: Uninstall command not found")
			return errors.New("uninstall command is empty")
		}

		commands := strings.Split(shellCommand, "&&")
		
		for _, cmd := range commands {
			parts := strings.Fields(strings.TrimSpace(cmd))
			if len(parts) == 0 {
				continue
			}
			execCmd := exec.Command(parts[0], parts[1:]...)
			err := execCmd.Run()
			if err != nil {
				log.Fatalf("Package Uninstallation Failed: %v", err)
				return err
			}
		}
	}

	fmt.Printf("Package %s uninstalled successfully.\n", packageName)
	return nil
}
