package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type registryPackage struct {
	Name         string   `json:"name"`
	Owner        string   `json:"owner"`
	Dependencies []string `json:"dependencies"`
}

func ListRegistryPackages(nameFilter, ownerFilter string) error {
	u, err := url.Parse("https://aetheis.vercel.app/packages")
	if err != nil {
		return err
	}
	q := u.Query()
	if nameFilter != "" {
		q.Set("name", nameFilter)
	}
	if ownerFilter != "" {
		q.Set("owner", ownerFilter)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to fetch packages: %w", err)
	}
	defer resp.Body.Close()

	if err := CheckAPIError(resp); err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var packages []registryPackage
	if err := json.Unmarshal(body, &packages); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Println("Packages on registry:")
	if len(packages) == 0 {
		fmt.Println("  (No packages)")
		return nil
	}
	for _, pkg := range packages {
		if len(pkg.Dependencies) > 0 {
			fmt.Printf("  - %s (owner: %s, deps: %v)\n", pkg.Name, pkg.Owner, pkg.Dependencies)
		} else {
			fmt.Printf("  - %s (owner: %s)\n", pkg.Name, pkg.Owner)
		}
	}
	return nil
}
