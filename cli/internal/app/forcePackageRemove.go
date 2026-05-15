package app

import (
	"errors"
	"fmt"
	"log"
)

// ForcePackageRemove forcefully removes a package from the store without dependency checks
// This bypasses the normal validation that prevents removal of packages that are required
// by other packages. It only removes the package from the store and does NOT attempt to
// uninstall it from the system.
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
		fmt.Printf("WARNING: Package %s is required by %s, but forcing removal from store anyway.\n", packageName, userPkg)
	}

	fmt.Printf("Force removing package %s from store (system files will remain)...\n", packageName)

	// Remove from Store only
	store.RemovePackage(packageName)
	if err := SavePackageStore(store); err != nil {
		log.Printf("Error saving store: %v", err)
		return err
	}
	fmt.Printf("Package %s has been removed from the package store.\n", packageName)

	return nil
}
