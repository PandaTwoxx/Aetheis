package app

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func Logout() error {
	fmt.Println("Logging out...")

	usr, err := user.Current()
	if err != nil {
		return err
	}

	tokenPath := filepath.Join(usr.HomeDir, ".aetheis", "token")
	os.Remove(tokenPath)
	fmt.Println("Successfully logged out.")
	return nil
}