package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MDMSResponse represents the structure of the response from the MDMS service.
type MDMSResponse struct {
	MDMS []struct {
		Data struct {
			NoOfDays int `json:"noOfDays"` // Number of days field in the MDMS data
		} `json:"data"`
	} `json:"mdms"`
}

// FetchNoOfDays fetches the number of days from the MDMS service for a given tenant.
// It sends a GET request to the MDMS endpoint, parses the response, and extracts the noOfDays field.
func FetchNoOfDays(url string, tenantID string) (int, error) {
	// Build the MDMS endpoint URL
	url = fmt.Sprintf("%s/mdms-v2/v2?schemaCode=PropertyTax.Enumeration.Duedate", url)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-ID", tenantID)
	httpReq.Header.Set("X-Client-Id", "test-client")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch MDMS data: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("MDMS service returned status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read MDMS response body: %w", err)
	}

	// Parse the response JSON into MDMSResponse struct
	var mdmsResp MDMSResponse
	if err := json.Unmarshal(body, &mdmsResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal MDMS response: %w", err)
	}

	// Extract the noOfDays field from the response
	if len(mdmsResp.MDMS) > 0 {
		return mdmsResp.MDMS[0].Data.NoOfDays, nil
	}

	return 0, fmt.Errorf("no data found in MDMS response")
}
