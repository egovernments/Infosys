// Package constants holds application-wide constant values
package constants

// Tenant and Workflow Constants
const (
	// Default tenant ID for the application
	DefaultTenantID = "pb.amritsar"

	// Workflow process type for property enumeration
	PropertyEnumerationWorkflow = "PROP_ENUM"

	// Application number formatting
	ApplicationNumberPrefix         = "APPL"
	ApplicationNumberDateFormat     = "20060102" // Date format for application number
	ApplicationNumberSequenceFormat = "%04d"     // Sequence format for application number
)

// User role constants
const (
	RoleCitizen        = "CITIZEN"
	RoleAgent          = "AGENT"
	RoleServiceManager = "SERVICE_MANAGER"
	RoleCommissioner   = "COMMISSIONER"
)

// Application status constants
const (
	StatusInitiated     = "INITIATED"
	StatusVerified      = "VERIFIED"
	StatusApproved      = "APPROVED"
	StatusAuditVerified = "AUDIT_VERIFIED"
	StatusAssigned      = "ASSIGNED"
	StatusRejected      = "REJECTED"
)

// Application priority constants
const (
	PriorityHigh   = "HIGH"
	PriorityMedium = "MEDIUM"
	PriorityLow    = "LOW"
	PriorityUrgent = "URGENT"
)

// Validation-related constants
const (
	MaxApplicationsPerCitizen = 10               // Max applications allowed per citizen
	MaxFileUploadSize         = 10 * 1024 * 1024 // 10MB file upload limit
	ApplicationNumberLength   = 20               // Length of application number
)

// Default values for pagination and application
const (
    DefaultPage = 0
    DefaultSize = 20
    MaxSize = 100
    
    // Add jurisdiction constants
    DefaultJurisdiction = "pb.amritsar"
    
    // Header names
    HeaderTenantID = "X-Tenant-ID"
    HeaderUserID = "X-User-ID"
    HeaderUserRole = "X-User-Role"
    HeaderStatus = "X-Status"
    HeaderTotalCount = "X-Total-Count"
    HeaderCurrentPage = "X-Current-Page"
    HeaderPerPage = "X-Per-Page"
    HeaderTotalPages = "X-Total-Pages"
)
