package app

import (
	"fmt"
	"io"
	"log"
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
		log.Fatalf("Error making GET request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error signing up: %v", resp.Status)
		return fmt.Errorf("error signing up: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return err
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
