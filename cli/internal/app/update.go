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

func Update(args []string) error {
	// 1. Load Store
	store, err := LoadPackageStore()
	if err != nil {
		log.Fatalf("Failed to load package store: %v", err)
		return err
	}
	if len(args) == 0 {
		// Update all packages
		targetPackage := InstalledPackage{
			Name:         "aetheis",
			Explicit:     true,
			Dependencies: []string{},
		}
		if store.GetPackage("aetheis") == nil {
			store.AddPackage(targetPackage)
			SavePackageStore(store)
		}
		for _, pkg := range store.ListPackages() {
			fmt.Printf("Updating package: %s\n", pkg)
			if pkg == "brew" || pkg == "brew-local" {
				currentUser, err := user.Current()
				if err != nil {
					log.Fatalf("Package Update Failed: %v", err)
					return err
				}
				cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
				_ = os.MkdirAll(cacheDir, 0755) // Ensure cache dir exists
				installPath := filepath.Join(cacheDir, "install_"+pkg+".sh")

				_ = os.Remove(installPath)
				// No need to Create then Open, just WriteFile or OpenFile with Create
				fileInstruct, err := os.OpenFile(installPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

				if err != nil {
					// Handle error
					log.Fatalf("Package Update Failed: %v", err)
					return err
				}
				fileInstruct.Write([]byte("#!/bin/sh\n" + "brew update" + "\n"))

				fileInstruct.Close()

				exec.Command("chmod", "+x", installPath).Run()

				execCmd := exec.Command("/bin/sh", installPath)
				// Connect stdout/stderr to see output
				execCmd.Stdout = os.Stdout
				execCmd.Stderr = os.Stderr

				cmdErr := execCmd.Run()
				if cmdErr != nil {
					log.Fatalf("Package Update Failed: %v", cmdErr)
					return cmdErr
				}
			} else if pkg == "aetheis" {
				fmt.Println("Updating aetheis...")
				if err := Update([]string{pkg}); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update aetheis: %v\n", err)
				}
			} else {
				if err := Update([]string{pkg}); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to update package %s: %v\n", pkg, err)
				}
			}
		}
	} else {
		for _, pkg := range args {
			fmt.Printf("Updating package: %s\n", pkg)
			if existing := store.GetPackage(pkg); existing == nil {
				fmt.Printf("Package %s is not installed.\n", pkg)
				continue
			}
			fmt.Println("Updating package:", pkg)
			resp, err := http.Get("https://aetheis.vercel.app/pkg/" + pkg)
			if err != nil {
				return fmt.Errorf("package installation failed: %w", err)
			}
			defer resp.Body.Close()

			if err := CheckAPIError(resp); err != nil {
				return err
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("package installation failed: %w", err)
			}

			parsed := strings.Split(string(bodyBytes), " ")

			sourceCode := strings.TrimSpace(parsed[0])

			if sourceCode == "" {
				return errors.New("package not found")
			}
			if sourceCode == "brew" {
				fmt.Printf("Installing via Homebrew: %s\n...", pkg)
				exec.Command("brew", "upgrade", pkg).Run()
			} else {
				currentUser, err := user.Current()
				if err != nil {
					log.Fatalf("Package Update Failed: %v", err)
					return err
				}
				// Try to fetch update script first
				resp, err := http.Get("https://aetheis.vercel.app/update/" + pkg)
				if err != nil {
					return fmt.Errorf("package update failed: %w", err)
				}
				defer resp.Body.Close()

				if err := CheckAPIError(resp); err != nil {
					return err
				}

				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("package update failed: %w", err)
				}

				shellScript := strings.TrimSpace(string(bodyBytes))
				if shellScript == "" {
					return errors.New("update command not found")
				}

				cacheDir := filepath.Join(currentUser.HomeDir, ".aetheis", "cache")
				_ = os.MkdirAll(cacheDir, 0755) // Ensure cache dir exists
				updatePath := filepath.Join(cacheDir, "update_"+pkg+".sh")

				_ = os.Remove(updatePath)

				// Write the full shell script to a file
				scriptContent := shellScript
				// Add shebang if not already present
				if !strings.HasPrefix(strings.TrimSpace(shellScript), "#!") {
					scriptContent = "#!/bin/sh\n" + shellScript
				}

				err = os.WriteFile(updatePath, []byte(scriptContent), 0755)
				if err != nil {
					log.Fatalf("Package Update Failed: %v", err)
					return err
				}

				// Execute the script
				execCmd := exec.Command("/bin/sh", updatePath)
				// Connect stdout/stderr to see output
				execCmd.Stdout = os.Stdout
				execCmd.Stderr = os.Stderr

				cmdErr := execCmd.Run()
				if cmdErr != nil {
					log.Fatalf("Package Update Failed: %v", cmdErr)
					return cmdErr
				}
			}
		}
	}
	return nil
}
