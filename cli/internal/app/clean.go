package app

import (
	"os"
	"os/user"
	"path/filepath"
)

// CleanCache removes all files in the ~/.aetheis/cache directory
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
			return nil
		}
		return err
	}

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

	return nil
}
