package app

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// CleanCache removes all files in the ~/.aetheis/cache directory
// and uninstalls unused dependencies of uninstalled packages
func CleanCache() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	cacheDir := filepath.Join(usr.HomeDir, ".aetheis", "cache")

	// Check if the cache directory exists
	if _, err := os.Stat(cacheDir); err != nil {
		if os.IsNotExist(err) {
			// Cache directory doesn't exist, nothing to clean
		} else {
			return err
		}
	} else {
		// Remove all files in the cache directory
		err = os.RemoveAll(cacheDir)
		if err != nil {
			return err
		}

		// Recreate the cache directory so it exists for future use
		err = os.MkdirAll(cacheDir, 0755)
		if err != nil {
			return err
		}
	}

	// Clean up unused dependencies
	if err := CleanUnusedDependencies(); err != nil {
		log.Printf("Warning: Failed to clean unused dependencies: %v\n", err)
	}

	return nil
}

// CleanUnusedDependencies identifies and uninstalls packages that are no longer needed
// (i.e., packages that were installed as dependencies but are no longer used by any other package)
func CleanUnusedDependencies() error {
	store, err := LoadPackageStore()
	if err != nil {
		return err
	}

	uninstalledCount := 0

	// Iterate through all packages
	for _, pkg := range store.Packages {
		// Skip explicitly installed packages (user-requested)
		if pkg.Explicit {
			continue
		}

		// Check if this implicit (dependency) package is used by any other package
		if _, used := store.IsPackageUsedByOthers(pkg.Name); !used {
			fmt.Printf("Removing unused dependency: %s\n", pkg.Name)

			// Uninstall the package
			if err := UninstallPackage(pkg.Name); err != nil {
				log.Printf("Warning: Failed to uninstall unused dependency %s: %v\n", pkg.Name, err)
				continue
			}
			uninstalledCount++

			// Reload the store since UninstallPackage modifies it
			store, err = LoadPackageStore()
			if err != nil {
				log.Printf("Warning: Failed to reload package store: %v\n", err)
				continue
			}
		}
	}

	if uninstalledCount > 0 {
		fmt.Printf("Cleaned up %d unused dependencies.\n", uninstalledCount)
	} else {
		fmt.Println("No unused dependencies found.")
	}

	return nil
}
