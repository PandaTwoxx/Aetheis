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

// Update is used to update aetheis itself using the registry
func Update(args []string) error {
	fmt.Println("🔄 Updating aetheis...")

	// Fetch aetheis update command from registry
	resp, err := http.Get("https://aetheis.vercel.app/update/aetheis")
	if err != nil {
		return fmt.Errorf("failed to fetch aetheis update command: %w", err)
	}
	defer resp.Body.Close()

	if err := CheckAPIError(resp); err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read aetheis update command: %w", err)
	}

	updateScript := strings.TrimSpace(string(bodyBytes))
	if updateScript == "" {
		return errors.New("aetheis update command not found in registry")
	}

	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
	_ = os.MkdirAll(cacheDir, 0755)
	updatePath := filepath.Join(cacheDir, "update_aetheis.sh")

	_ = os.Remove(updatePath)

	// Write the update script to a file
	scriptContent := updateScript
	if !strings.HasPrefix(strings.TrimSpace(updateScript), "#!") {
		scriptContent = "#!/bin/sh\n" + updateScript
	}

	err = os.WriteFile(updatePath, []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to write update script: %w", err)
	}

	fmt.Println("📝 Executing aetheis update script...")
	execCmd := exec.Command("/bin/sh", updatePath)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("aetheis update failed: %w", err)
	}

	fmt.Println("✅ Aetheis updated successfully!")
	return nil
}

// Upgrade updates installed packages (what Update used to do)
func Upgrade(args []string) error {
	// 1. Load Store
	store, err := LoadPackageStore()
	if err != nil {
		log.Fatalf("Failed to load package store: %v", err)
		return err
	}
	if len(args) == 0 {
		// Upgrade all packages
		for _, pkg := range store.ListPackages() {
			fmt.Printf("Upgrading package: %s\n", pkg)
			if pkg == "brew" || pkg == "brew-local" {
				currentUser, err := user.Current()
				if err != nil {
					log.Fatalf("Package Upgrade Failed: %v", err)
					return err
				}
				cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
				_ = os.MkdirAll(cacheDir, 0755) // Ensure cache dir exists
				upgradePath := filepath.Join(cacheDir, "upgrade_"+pkg+".sh")

				_ = os.Remove(upgradePath)
				fileInstruct, err := os.OpenFile(upgradePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

				if err != nil {
					log.Fatalf("Package Upgrade Failed: %v", err)
					return err
				}
				fileInstruct.Write([]byte("#!/bin/sh\n" + "brew update" + "\n"))
				fileInstruct.Close()

				exec.Command("chmod", "+x", upgradePath).Run()

				execCmd := exec.Command("/bin/sh", upgradePath)
				execCmd.Stdout = os.Stdout
				execCmd.Stderr = os.Stderr

				cmdErr := execCmd.Run()
				if cmdErr != nil {
					log.Fatalf("Package Upgrade Failed: %v", cmdErr)
					return cmdErr
				}
			} else if pkg != "aetheis" {
				if err := Upgrade([]string{pkg}); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to upgrade package %s: %v\n", pkg, err)
				}
			}
		}
	} else {
		for _, pkg := range args {
			fmt.Printf("Upgrading package: %s\n", pkg)
			if existing := store.GetPackage(pkg); existing == nil {
				fmt.Printf("Package %s is not installed.\n", pkg)
				continue
			}
			fmt.Println("Upgrading package:", pkg)
			resp, err := http.Get("https://aetheis.vercel.app/pkg/" + pkg)
			if err != nil {
				return fmt.Errorf("package upgrade failed: %w", err)
			}
			defer resp.Body.Close()

			if err := CheckAPIError(resp); err != nil {
				return err
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("package upgrade failed: %w", err)
			}

			parsed := strings.Split(string(bodyBytes), " ")
			sourceCode := strings.TrimSpace(parsed[0])

			if sourceCode == "" {
				return errors.New("package not found")
			}
			if sourceCode == "brew" {
				fmt.Printf("Upgrading via Homebrew: %s\n...", pkg)
				exec.Command("brew", "upgrade", pkg).Run()
			} else {
				currentUser, err := user.Current()
				if err != nil {
					log.Fatalf("Package Upgrade Failed: %v", err)
					return err
				}
				// Try to fetch update script first
				resp, err := http.Get("https://aetheis.vercel.app/update/" + pkg)
				if err != nil {
					return fmt.Errorf("package upgrade failed: %w", err)
				}
				defer resp.Body.Close()

				if err := CheckAPIError(resp); err != nil {
					return err
				}

				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("package upgrade failed: %w", err)
				}

				shellScript := strings.TrimSpace(string(bodyBytes))
				if shellScript == "" {
					return errors.New("update command not found")
				}

				cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
				_ = os.MkdirAll(cacheDir, 0755) // Ensure cache dir exists
				upgradePath := filepath.Join(cacheDir, "upgrade_"+pkg+".sh")

				_ = os.Remove(upgradePath)

				// Write the full shell script to a file
				scriptContent := shellScript
				// Add shebang if not already present
				if !strings.HasPrefix(strings.TrimSpace(shellScript), "#!") {
					scriptContent = "#!/bin/sh\n" + shellScript
				}

				err = os.WriteFile(upgradePath, []byte(scriptContent), 0755)
				if err != nil {
					log.Fatalf("Package Upgrade Failed: %v", err)
					return err
				}

				// Execute the script
				execCmd := exec.Command("/bin/sh", upgradePath)
				execCmd.Stdout = os.Stdout
				execCmd.Stderr = os.Stderr

				cmdErr := execCmd.Run()
				if cmdErr != nil {
					log.Fatalf("Package Upgrade Failed: %v", cmdErr)
					return cmdErr
				}
			}
		}
	}
	return nil
}
