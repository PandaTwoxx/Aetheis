package app

import (
	"fmt"
	"os"
	"net/http"
	"os/user"
	"path/filepath"
)

func AddPackage(packageName string) error {
	installScript := ""
	uninstallScript := ""

	fmt.Println("Enter the install script for the package:")
	fmt.Scan(&installScript)
	fmt.Println("Enter the uninstall script for the package:")
	fmt.Scan(&uninstallScript)
	
	dependencyList := ""

	fmt.Println("Enter the dependencies for the package (separated by spaces):")
	fmt.Scan(&dependencyList)

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


	uploadLink := "https://aetheis.vercel.app/addPackage/" + token + "/" + packageName + "/" + installScript + "/" + uninstallScript + "/" + dependencyList

	resp, err := http.Post(uploadLink, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Package added successfully.")
	
	return nil
}