package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func Signup(username string, password string) error {
	fmt.Println("Signing up...")

	url := "https://aetheis.vercel.app/addUser/" + username + "/" + password

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if err := CheckAPIError(resp); err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	tokenPath := filepath.Join(usr.HomeDir, ".aetheis", "token")

	os.Remove(tokenPath)
	os.Create(tokenPath)

	file, err := os.OpenFile(tokenPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	file.Write(body)
	file.Close()

	fmt.Println("Successfully logged in.")

	return nil
}
