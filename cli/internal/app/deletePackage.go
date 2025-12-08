package app

import (
	"fmt"
	"os"
	"net/http"
	"os/user"
	"path/filepath"
)

func DeletePackage(packageName string) error {

	user, err := user.Current()
	if err != nil {
		return err
	}

	path := filepath.Join(user.HomeDir, ".aetheis", "token")
	
	tokenFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	token := string(tokenFile)


	uploadLink := "https://aetheis.vercel.app/deletePackage/" + token + "/" + packageName

	resp, err := http.Post(uploadLink, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Package deleted successfully.")
	
	return nil
}