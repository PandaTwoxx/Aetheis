package app

import (
	"fmt"
	"log"
	"errors"
)

func ForcePackageAdd(packageName string) error {
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
		existing.Explicit = true
		SavePackageStore(store)
		fmt.Printf("Package %s is already installed.\n", packageName)
		// We could return nil here, but maybe we want to verify dependencies again?
		// Let's assume re-install is okay or we just return.
		// For matching original behavior regarding dependencies, we might want to check them.
		// But let's proceed with install to be safe.
	}

	if packageName == "" {
		return errors.New("package name cannot be empty")
	}

	targetPackage := InstalledPackage{
		Name:         packageName,
		Explicit:     true,
		Dependencies: []string{},
	}

	store.AddPackage(targetPackage)
	if err := SavePackageStore(store); err != nil {
		log.Printf("Warning: Failed to save package store: %v", err)
	}

	fmt.Printf("Package %s added successfully.\n", packageName)

	return nil
}