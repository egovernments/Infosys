package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// NotificationRequest represents the payload for sending a notification.
type NotificationRequest struct {
	Title  string                 `json:"title"`  // Notification title
	Body   string                 `json:"body"`   // Notification body/message
	Data   map[string]interface{} `json:"data"`   // Additional data to send with the notification
	Tokens []string               `json:"tokens"` // Device tokens to which the notification will be sent
}

// SendNotification sends a notification to the specified URL using the provided NotificationRequest payload.
// It marshals the request to JSON, sends a POST request, and checks for a successful response.
func SendNotification(url string, req NotificationRequest) error {
	// Marshal the notification request to JSON
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal notification request: %w", err)
	}

	// Create the HTTP POST request
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification service returned status: %s", resp.Status)
	}

	return nil
}
