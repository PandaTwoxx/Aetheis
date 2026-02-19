package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func UpdatePackage(packageName string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the install script for the package:")
	installScript, _ := reader.ReadString('\n')
	installScript = strings.TrimSpace(installScript)

	fmt.Println("Enter the uninstall script for the package:")
	uninstallScript, _ := reader.ReadString('\n')
	uninstallScript = strings.TrimSpace(uninstallScript)

	fmt.Println("Enter the update script for the package:")
	updateScript, _ := reader.ReadString('\n')
	updateScript = strings.TrimSpace(updateScript)

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
	token := strings.TrimSpace(string(tokenFile))

	dependencies := []string{}
	if dependencyList != "" {
		dependencies = strings.Split(dependencyList, " ")
	}

	body, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"name":              packageName,
		"installCommands":   installScript,
		"uninstallCommands": uninstallScript,
		"updateCommands":    updateScript,
		"dependencies":      dependencies,
	})
	if err != nil {
		return err
	}

	uploadLink := "https://aetheis.vercel.app/updatePackage"
	resp, err := http.Post(uploadLink, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := CheckAPIError(resp); err != nil {
		return err
	}

	fmt.Println("Package updated successfully.")
	
	return nil
}