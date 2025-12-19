package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UserResponse represents the structure of the response from the user service API.
type UserResponse struct {
	Success bool   `json:"success"` // Indicates if the API call was successful
	Message string `json:"message"` // Message from the API (error or success info)
	Data    struct {
		ID       string `json:"id"`       // User ID
		Username string `json:"username"` // Username of the user
		// Other fields omitted for brevity
	} `json:"data"`
}

// GetUsername fetches the username for a given userID from the user service.
// It sends a GET request to the user service API, decodes the response, and returns the username.
func GetUsername(userID string) (string, error) {
	// Build the user service API URL
	url := fmt.Sprintf("http://10.232.161.103:30105/api/v1/users/%s", userID)

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into UserResponse struct
	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return "", err
	}

	// Check if the API call was successful
	if !userResp.Success {
		return "", fmt.Errorf("API error: %s", userResp.Message)
	}

	// Return the username from the response
	return userResp.Data.Username, nil
}
