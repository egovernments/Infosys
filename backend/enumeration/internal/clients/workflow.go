package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is a workflow service client for making HTTP requests to the workflow API.
type Client struct {
	baseURL    string       // Base URL of the workflow service
	httpClient *http.Client // HTTP client used for requests
}

// Process represents a workflow process definition.
type Process struct {
	ID   string `json:"id"`   // Unique process ID
	Code string `json:"code"` // Process code
	Name string `json:"name"` // Process name
}

// ProcessResponse represents the response from a workflow transition or process action.
type ProcessResponse struct {
	ID           string              `json:"id"`           // Unique response ID
	ProcessID    string              `json:"processId"`    // Associated process ID
	EntityID     string              `json:"entityId"`     // Entity ID involved in the process
	Action       string              `json:"action"`       // Action performed
	Status       string              `json:"status"`       // Status of the process
	Comment      string              `json:"comment"`      // Comment for the action
	Documents    []string            `json:"documents"`    // Related documents
	Assigner     string              `json:"assigner"`     // User who assigned the task
	Assignees    []string            `json:"assignees"`    // Users assigned to the task
	CurrentState string              `json:"currentState"` // Current state of the process
	StateSla     int64               `json:"stateSla"`     // State SLA in milliseconds
	ProcessSla   int64               `json:"processSla"`   // Process SLA in milliseconds
	Attributes   map[string][]string `json:"attributes"`   // Additional attributes
	AuditDetails AuditDetail         `json:"auditDetails"` // Audit details
}

// AuditDetail contains audit information for workflow entities.
type AuditDetail struct {
	CreatedBy    string `json:"createdBy,omitempty" db:"created_by" gorm:"column:created_by"`      // User who created the entity
	CreatedTime  int64  `json:"createdTime,omitempty" db:"created_at" gorm:"column:created_at"`    // Creation timestamp
	ModifiedBy   string `json:"modifiedBy,omitempty" db:"modified_by" gorm:"column:modified_by"`   // User who last modified the entity
	ModifiedTime int64  `json:"modifiedTime,omitempty" db:"modified_at" gorm:"column:modified_at"` // Last modification timestamp
}

// StateResponse represents a workflow state.
type StateResponse struct {
	ID          string `json:"id"`          // State ID
	Code        string `json:"code"`        // State code
	Name        string `json:"name"`        // State name
	Description string `json:"description"` // State description
}

// ActionResponse represents a workflow action and its next state.
type ActionResponse struct {
	ID        string `json:"id"`        // Action ID
	Name      string `json:"name"`      // Action name
	NextState string `json:"nextState"` // Next state after action
}

// TransitionRequest represents a request to transition a workflow process.
type TransitionRequest struct {
	ProcessID  string                 `json:"processId"`  // ID of the process to transition
	EntityID   string                 `json:"entityId"`   // ID of the entity involved
	Action     string                 `json:"action"`     // Action to perform
	Attributes map[string]interface{} `json:"attributes"` // Additional attributes for the transition
}

// NewClient creates a new workflow Client with the given base URL and a default HTTP client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetProcess fetches a workflow process definition by process code for a given tenant.
// Returns the process if found, or an error if not found or on failure.
func (c *Client) GetProcess(ctx context.Context, tenantID, processCode string) (*Process, error) {
	url := fmt.Sprintf("%s/workflow/v1/process?code=%s", c.baseURL, processCode)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Tenant-ID", tenantID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workflow service returned status %d", resp.StatusCode)
	}

	var processes []Process
	if err := json.NewDecoder(resp.Body).Decode(&processes); err != nil {
		return nil, err
	}

	for _, p := range processes {
		if p.Code == processCode {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("process with code %s not found", processCode)
}

// CreateTransition sends a transition request to the workflow service to perform an action on a process.
// It builds the payload, sends a POST request, and returns the resulting process response.
func (c *Client) CreateTransition(ctx context.Context, tenantID string, req *TransitionRequest) (*ProcessResponse, error) {
	url := fmt.Sprintf("%s/workflow/v1/transition", c.baseURL)

	// Add comment field as per workflow service requirement
	transitionPayload := map[string]interface{}{
		"processId":  req.ProcessID,
		"entityId":   req.EntityID,
		"action":     req.Action,
		"comment":    fmt.Sprintf("Executing %s action", req.Action),
		"attributes": req.Attributes,
	}

	jsonData, err := json.Marshal(transitionPayload)
	if err != nil {
		return nil, err
	}

	// Debug: print the request payload
	fmt.Printf("DEBUG: Transition request: %s\n", string(jsonData))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-ID", tenantID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("DEBUG: Workflow service returned status %d for URL: %s\n", resp.StatusCode, url)
		return nil, fmt.Errorf("workflow service returned status %d", resp.StatusCode)
	}

	var result ProcessResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
