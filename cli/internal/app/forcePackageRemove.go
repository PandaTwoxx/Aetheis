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

// ForcePackageRemove forcefully removes a package without dependency checks
// This bypasses the normal validation that prevents removal of packages
// that are required by other packages.
func ForcePackageRemove(packageName string) error {
	store, err := LoadPackageStore()
	if err != nil {
		log.Fatalf("Failed to load package store: %v", err)
		return err
	}

	pkg := store.GetPackage(packageName)
	if pkg == nil {
		log.Printf("Package %s not found in package store.", packageName)
		return errors.New("package not found")
	}

	// Check if package is used by others and warn
	if userPkg, used := store.IsPackageUsedByOthers(packageName); used {
		fmt.Printf("WARNING: Package %s is required by %s, but forcing removal anyway.\n", packageName, userPkg)
	}

	fmt.Printf("Force removing package: %s\n", packageName)

	// Perform Uninstall
	resp, err := http.Get("https://aetheis.vercel.app/pkg/" + packageName)
	if err != nil {
		log.Printf("Warning: Failed to fetch package info for uninstall: %v. Proceeding with removal from store.", err)
	} else {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Printf("Package definition not found on server. Proceeding with removal from local store.\n")
		} else if err := CheckAPIError(resp); err != nil {
			log.Printf("Warning: Failed to check API error: %v. Proceeding with removal from store.", err)
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			sourceInfo := strings.TrimSpace(strings.Split(string(bodyBytes), " ")[0])

			if sourceInfo == "brew" {
				fmt.Printf("Uninstalling via Homebrew: %s\n...", packageName)
				exec.Command("brew", "uninstall", packageName).Run()
			} else {
				fmt.Printf("Uninstalling via Shell Script\n...")
				// Fetch uninstall script
				resp, err := http.Get("https://aetheis.vercel.app/uninstall/" + packageName)
				if err != nil {
					log.Printf("Warning: Failed to fetch uninstall script: %v. Removing from store anyway.", err)
				} else {
					defer resp.Body.Close()
					if err := CheckAPIError(resp); err != nil {
						log.Printf("Warning: Failed to check API error: %v. Removing from store anyway.", err)
					} else {
						uninstallBytes, _ := io.ReadAll(resp.Body)
						shellScript := string(uninstallBytes)

						if shellScript != "" {
							// Create a temporary file for the uninstall script
							currentUser, err := user.Current()
							if err != nil {
								log.Printf("Warning: Could not get current user: %v. Running script directly.", err)
								execCmd := exec.Command("sh", "-c", shellScript)
								execCmd.Stdout = os.Stdout
								execCmd.Stderr = os.Stderr
								if err := execCmd.Run(); err != nil {
									log.Printf("Warning: Uninstall script failed: %v. Removing from store anyway.", err)
								}
							} else {
								cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
								_ = os.MkdirAll(cacheDir, 0755)
								uninstallPath := filepath.Join(cacheDir, "uninstall_"+packageName+".sh")
								_ = os.Remove(uninstallPath)

								// Write the full shell script to a file
								scriptContent := shellScript
								// Add shebang if not already present
								if !strings.HasPrefix(strings.TrimSpace(shellScript), "#!") {
									scriptContent = "#!/bin/sh\n" + shellScript
								}

								err = os.WriteFile(uninstallPath, []byte(scriptContent), 0755)
								if err != nil {
									log.Printf("Warning: Failed to write uninstall script: %v. Running inline.", err)
									execCmd := exec.Command("sh", "-c", shellScript)
									execCmd.Stdout = os.Stdout
									execCmd.Stderr = os.Stderr
									if err := execCmd.Run(); err != nil {
										log.Printf("Warning: Uninstall script failed: %v. Removing from store anyway.", err)
									}
								} else {
									execCmd := exec.Command("/bin/sh", uninstallPath)
									execCmd.Stdout = os.Stdout
									execCmd.Stderr = os.Stderr
									if err := execCmd.Run(); err != nil {
										log.Printf("Warning: Uninstall script failed: %v. Removing from store anyway.", err)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Remove from Store
	store.RemovePackage(packageName)
	if err := SavePackageStore(store); err != nil {
		log.Printf("Error saving store: %v", err)
	}
	fmt.Printf("Package %s force removed successfully.\n", packageName)

	return nil
}
