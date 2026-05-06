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

	fmt.Println("Enter the path to the install script:")
	installScriptPath, _ := reader.ReadString('\n')
	installScriptPath = strings.TrimSpace(installScriptPath)

	installScript, err := readScriptFile(installScriptPath)
	if err != nil {
		return fmt.Errorf("failed to read install script: %w", err)
	}

	fmt.Println("Enter the path to the uninstall script:")
	uninstallScriptPath, _ := reader.ReadString('\n')
	uninstallScriptPath = strings.TrimSpace(uninstallScriptPath)

	uninstallScript, err := readScriptFile(uninstallScriptPath)
	if err != nil {
		return fmt.Errorf("failed to read uninstall script: %w", err)
	}

	fmt.Println("Enter the path to the update script:")
	updateScriptPath, _ := reader.ReadString('\n')
	updateScriptPath = strings.TrimSpace(updateScriptPath)

	updateScript, err := readScriptFile(updateScriptPath)
	if err != nil {
		return fmt.Errorf("failed to read update script: %w", err)
	}

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