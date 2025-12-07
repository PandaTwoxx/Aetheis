package app

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type InstalledPackage struct {
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies"`
	Explicit     bool     `json:"explicit"` // True if installed directly by user, False if dependency
}

type PackageStore struct {
	Packages []InstalledPackage `json:"packages"`
}

func GetStorePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".aetheis", "install_packages.json"), nil
}

func LoadPackageStore() (*PackageStore, error) {
	path, err := GetStorePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &PackageStore{Packages: []InstalledPackage{}}, nil
	}
	if err != nil {
		return nil, err
	}

	// Handle empty file (0 bytes) or empty string
	if len(data) == 0 {
		return &PackageStore{Packages: []InstalledPackage{}}, nil
	}

	// Try parsing as JSON first
	var store PackageStore
	// Check if it's the old format (newline separated list) or empty
	if len(data) > 0 && data[0] != '{' && data[0] != '[' {
		// Old format migration
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				store.Packages = append(store.Packages, InstalledPackage{
					Name:     trimmed,
					Explicit: true, // Assume explicit for migrated legacy packages
				})
			}
		}
		// Save immediately to migrate
		_ = SavePackageStore(&store)
		return &store, nil
	}

	// Support if the root is just the array
	if len(data) > 0 && data[0] == '[' {
		var pkgs []InstalledPackage
		if err := json.Unmarshal(data, &pkgs); err != nil {
			return nil, err
		}
		store.Packages = pkgs
		return &store, nil
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

func SavePackageStore(store *PackageStore) error {
	path, err := GetStorePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *PackageStore) AddPackage(pkg InstalledPackage) {
	for i, p := range s.Packages {
		if p.Name == pkg.Name {
			// Update existing
			s.Packages[i] = pkg
			return
		}
	}
	s.Packages = append(s.Packages, pkg)
}

func (s *PackageStore) RemovePackage(name string) {
	var newPkgs []InstalledPackage
	for _, p := range s.Packages {
		if p.Name != name {
			newPkgs = append(newPkgs, p)
		}
	}
	s.Packages = newPkgs
}

func (s *PackageStore) GetPackage(name string) *InstalledPackage {
	for _, p := range s.Packages {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

// IsPackageUsedByOthers returns true if 'name' is listed as a dependency of any OTHER package
func (s *PackageStore) IsPackageUsedByOthers(name string) (string, bool) {
	for _, p := range s.Packages {
		if p.Name == name {
			continue
		}
		for _, dep := range p.Dependencies {
			if dep == name {
				return p.Name, true
			}
		}
	}
	return "", false
}
