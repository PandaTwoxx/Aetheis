package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CheckAPIError reads the response body and returns an error if the status code
// indicates failure (4xx or 5xx). It attempts to parse JSON error responses
// with an "err" field from the server.
func CheckAPIError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("server error (%d): failed to read response", resp.StatusCode)
	}
	var result struct {
		Err string `json:"err"`
	}
	if json.Unmarshal(body, &result) == nil && result.Err != "" {
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, result.Err)
	}
	return fmt.Errorf("server error (%d): %s", resp.StatusCode, string(body))
}
