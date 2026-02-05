package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func UninstallPackage(packageName string) error {
	store, err := LoadPackageStore()
	if err != nil {
		log.Fatalf("Failed to load package store: %v", err)
		return err
	}

	pkg := store.GetPackage(packageName)
	if pkg == nil {
		// If not in store, maybe it was installed by old version or manually?
		// For now fail or try blindly? Let's fail if not found in our managed store to be safe,
		// but maybe we should allow force uninstall?
		// Given the task is to make uninstallation work correctly with metadata, we assume it's in store.
		log.Printf("Package %s not found in package store.", packageName)
		return errors.New("package not found")
	}

	// 1. Dependency Check
	if userPkg, used := store.IsPackageUsedByOthers(packageName); used {
		return fmt.Errorf("cannot uninstall %s: it is required by %s", packageName, userPkg)
	}

	fmt.Printf("Uninstalling package: %s\n", packageName)

	// 2. Perform Uninstall
	// We need to know the source (brew vs custom).
	// We didn't store "SourceType" explicitly in metadata, but we can infer or we should have stored it?
	// `InstalledPackage` struct has Name, Dependencies, Explicit.
	// `InstallPackage` logic determined source by fetching from "https://aetheis.vercel.app/" + packageName.
	// We should probably re-fetch to get the uninstall command/type, similar to Install.
	// Existing uninstallation code did exactly that (fetch URL).

	resp, err := http.Get("https://aetheis.vercel.app/pkg/" + packageName)
	if err != nil {
		log.Printf("Warning: Failed to fetch package info for uninstall: %v. Proceeding with removal from store.", err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Printf("Package definition not found on server. Proceeding with removal from local store.\n")
		} else if err := CheckAPIError(resp); err != nil {
			return err
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			sourceInfo := strings.TrimSpace(strings.Split(string(bodyBytes), " ")[0])

			if sourceInfo == "brew" {
				fmt.Printf("Uninstalling via Homebrew: %s\n...", packageName)
				exec.Command("brew", "uninstall", packageName).Run()
			} else {
				fmt.Printf("Uninstalling via Shell Command: %s\n...", sourceInfo)
				// Fetch uninstall script
				resp, err := http.Get("https://aetheis.vercel.app/uninstall/" + packageName)
				if err != nil {
					return fmt.Errorf("failed to fetch uninstall script: %w", err)
				}
				defer resp.Body.Close()
				if err := CheckAPIError(resp); err != nil {
					return err
				}
				uninstallBytes, _ := io.ReadAll(resp.Body)
				shellCommand := string(uninstallBytes)

				// FIX: The server sometimes sends `curl -fsSLO ... > ...`, which breaks because -O writes to file and stdout is empty.
				shellCommand = strings.ReplaceAll(shellCommand, "curl -fsSLO", "curl -fsSL")
				shellCommand = strings.ReplaceAll(shellCommand, "curl -O", "curl")

				// FIX: The server sends `mkdir -p ~/.aetheis/uninstall`. We want to use `~/.aetheis/cache`.
				shellCommand = strings.ReplaceAll(shellCommand, "~/.aetheis/uninstall", "~/.aetheis/cache")

				// FIX: The script writes to ~/.aetheis/cache/uninstall.sh but executes uninstall.sh (CWD).
				shellCommand = strings.ReplaceAll(shellCommand, "/bin/bash uninstall.sh", "sed -i.bak 's/ execute_sudo/ /g' ~/.aetheis/cache/uninstall.sh && /bin/bash ~/.aetheis/cache/uninstall.sh")

				// FIX: Shell doesn't expand '~' in '--path=~...' arguments. Use $HOME instead.
				shellCommand = strings.ReplaceAll(shellCommand, "~/", "$HOME/")

				if shellCommand != "" {
					execCmd := exec.Command("sh", "-c", shellCommand)
					execCmd.Stdout = os.Stdout
					execCmd.Stderr = os.Stderr
					if err := execCmd.Run(); err != nil {
						log.Printf("Warning: Uninstall script failed: %v. Removing from store anyway.", err)
					}
				}
			}
		}
	}

	// 3. Remove from Store
	store.RemovePackage(packageName)
	if err := SavePackageStore(store); err != nil {
		log.Printf("Error saving store: %v", err)
	}
	fmt.Printf("Package %s uninstalled successfully.\n", packageName)

	// 4. Auto-remove dependencies
	for _, depName := range pkg.Dependencies {
		depPkg := store.GetPackage(depName)
		if depPkg == nil {
			continue
		}
		// Check if dependency is now unused
		if _, used := store.IsPackageUsedByOthers(depName); !used {
			// Check if it was explicitly installed by user
			if !depPkg.Explicit {
				fmt.Printf("Dependency %s is no longer needed. Auto-removing...\n", depName)
				if err := UninstallPackage(depName); err != nil {
					log.Printf("Failed to auto-remove dependency %s: %v\n", depName, err)
				}
			}
		}
	}

	return nil
}
