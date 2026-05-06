package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// PreviewInstallPackage shows what commands will be executed during installation without actually running them
func PreviewInstallPackage(packageName string) error {
	if packageName == "" {
		return errors.New("package name cannot be empty")
	}

	store, err := LoadPackageStore()
	if err != nil {
		log.Printf("Warning: Could not load package store: %v", err)
	}

	// Check if already installed
	if store != nil && store.GetPackage(packageName) != nil {
		fmt.Printf("⚠️  Package %s is already installed.\n", packageName)
		return nil
	}

	fmt.Printf("\n📦 Package: %s\n", packageName)
	fmt.Println(strings.Repeat("-", 50))

	// Fetch package information
	resp, err := http.Get("https://aetheis.vercel.app/pkg/" + packageName)
	if err != nil {
		return fmt.Errorf("failed to fetch package info: %w", err)
	}
	defer resp.Body.Close()

	if err := CheckAPIError(resp); err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read package info: %w", err)
	}

	parsed := strings.Split(string(bodyBytes), " ")
	sourceCode := strings.TrimSpace(parsed[0])
	var dependencies []string
	if len(parsed) > 1 {
		dependencies = parsed[1:]
	}

	if sourceCode == "" {
		return errors.New("package not found")
	}

	// Show installation method
	fmt.Println("\n📋 Installation Method:")
	if sourceCode == "brew" {
		fmt.Println("  Command: brew install " + packageName)
	} else {
		fmt.Println("  Method: Custom shell script")
		fmt.Println("  Endpoint: https://aetheis.vercel.app/install/" + packageName)

		// Fetch and display the install script
		resp, err := http.Get("https://aetheis.vercel.app/install/" + packageName)
		if err == nil {
			defer resp.Body.Close()
			scriptBytes, _ := io.ReadAll(resp.Body)
			scriptContent := strings.TrimSpace(string(scriptBytes))
			if scriptContent != "" {
				fmt.Println("\n  Script preview (first 500 chars):")
				preview := scriptContent
				if len(preview) > 500 {
					preview = preview[:500] + "...[truncated]"
				}
				// Indent the script
				for _, line := range strings.Split(preview, "\n") {
					fmt.Printf("    %s\n", line)
				}
			}
		}
	}

	// Show dependencies
	if len(dependencies) > 0 {
		fmt.Println("\n📌 Dependencies (will be installed first):")
		for _, dep := range dependencies {
			fmt.Printf("  - %s\n", strings.TrimSpace(dep))
		}
	} else {
		fmt.Println("\n📌 Dependencies: None")
	}

	fmt.Println("\n✅ This preview shows what will be executed during installation.")
	fmt.Println("   Run 'aetheis install " + packageName + "' to proceed.\n")

	return nil
}

// PreviewUninstallPackage shows what commands will be executed during uninstallation
func PreviewUninstallPackage(packageName string) error {
	if packageName == "" {
		return errors.New("package name cannot be empty")
	}

	store, err := LoadPackageStore()
	if err != nil {
		return fmt.Errorf("failed to load package store: %w", err)
	}

	pkg := store.GetPackage(packageName)
	if pkg == nil {
		return errors.New("package not found in store")
	}

	fmt.Printf("\n📦 Package: %s\n", packageName)
	fmt.Println(strings.Repeat("-", 50))

	// Check if package is used by others
	if userPkg, used := store.IsPackageUsedByOthers(packageName); used {
		fmt.Printf("\n⚠️  WARNING: Package %s is required by %s\n", packageName, userPkg)
		fmt.Println("   Uninstallation would break this dependency!")
		return nil
	}

	// Fetch package information
	resp, err := http.Get("https://aetheis.vercel.app/pkg/" + packageName)
	if err != nil {
		fmt.Println("\n⚠️  Warning: Could not fetch package info from server")
		fmt.Println("  But it will be removed from the local store anyway")
		goto showDependents
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("\n⚠️  Package definition not found on server")
		fmt.Println("  It will be removed from the local store")
	} else if err := CheckAPIError(resp); err == nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		sourceInfo := strings.TrimSpace(strings.Split(string(bodyBytes), " ")[0])

		fmt.Println("\n📋 Uninstallation Method:")
		if sourceInfo == "brew" {
			fmt.Println("  Command: brew uninstall " + packageName)
		} else {
			fmt.Println("  Method: Custom shell script")
			fmt.Println("  Endpoint: https://aetheis.vercel.app/uninstall/" + packageName)

			// Fetch and display the uninstall script
			resp, err := http.Get("https://aetheis.vercel.app/uninstall/" + packageName)
			if err == nil {
				defer resp.Body.Close()
				scriptBytes, _ := io.ReadAll(resp.Body)
				scriptContent := strings.TrimSpace(string(scriptBytes))
				if scriptContent != "" {
					fmt.Println("\n  Script preview (first 500 chars):")
					preview := scriptContent
					if len(preview) > 500 {
						preview = preview[:500] + "...[truncated]"
					}
					// Indent the script
					for _, line := range strings.Split(preview, "\n") {
						fmt.Printf("    %s\n", line)
					}
				}
			}
		}
	}

showDependents:
	// Show what depends on this package
	if len(pkg.Dependencies) > 0 {
		fmt.Println("\n📌 Packages that depend on this:")
		for _, dep := range pkg.Dependencies {
			fmt.Printf("  - %s\n", dep)
		}
	}

	fmt.Println("\n✅ This preview shows what will be executed during uninstallation.")
	fmt.Println("   Run 'aetheis uninstall " + packageName + "' to proceed.\n")

	return nil
}

// PreviewUpdate shows what commands will be executed during update
func PreviewUpdate(packageList []string) error {
	store, err := LoadPackageStore()
	if err != nil {
		return fmt.Errorf("failed to load package store: %w", err)
	}

	if len(packageList) == 0 {
		// Update all packages
		fmt.Println("\n🔄 Updating all installed packages:\n")
		installedPackages := store.ListPackages()

		if len(installedPackages) == 0 {
			fmt.Println("No packages installed.")
			return nil
		}

		fmt.Printf("Packages to update (%d total):\n", len(installedPackages))
		for _, pkg := range installedPackages {
			if pkg == "brew" || pkg == "brew-local" {
				fmt.Printf("  • %s -> brew update\n", pkg)
			} else {
				fmt.Printf("  • %s -> Check for updates\n", pkg)
			}
		}
	} else {
		fmt.Printf("\n🔄 Updating packages: %v\n", packageList)
		fmt.Println(strings.Repeat("-", 50))

		for _, pkg := range packageList {
			fmt.Printf("\n📦 Package: %s\n", pkg)

			if store.GetPackage(pkg) == nil {
				fmt.Printf("  ⚠️  Not installed\n")
				continue
			}

			if pkg == "brew" || pkg == "brew-local" {
				fmt.Println("  Command: brew update")
			} else {
				fmt.Println("  Method: Fetch latest version and update")
				fmt.Printf("  Endpoint: https://aetheis.vercel.app/pkg/%s\n", pkg)
			}
		}
	}

	fmt.Println("\n✅ This preview shows what will be updated.")
	if len(packageList) > 0 {
		fmt.Printf("   Run 'aetheis update %s' to proceed.\n", strings.Join(packageList, " "))
	} else {
		fmt.Println("   Run 'aetheis update' to proceed.\n")
	}

	return nil
}
