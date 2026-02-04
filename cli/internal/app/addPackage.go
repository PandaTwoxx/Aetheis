package app

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func AddPackage(packageName string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the install script for the package:")
	installScript, _ := reader.ReadString('\n')
	installScript = strings.TrimSpace(installScript)

	fmt.Println("Enter the uninstall script for the package:")
	uninstallScript, _ := reader.ReadString('\n')
	uninstallScript = strings.TrimSpace(uninstallScript)

	fmt.Println("Enter the dependencies for the package (separated by spaces):")
	dependencyList, _ := reader.ReadString('\n')
	dependencyList = strings.TrimSpace(dependencyList)

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