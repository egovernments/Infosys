// Package dto defines data transfer objects for API requests and responses
package dto

import (
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// APPLICATION DTOs
// ============================================================================

// CreateApplicationRequest is used to create a new application
type CreateApplicationRequest struct {
	PropertyID    uuid.UUID `json:"propertyId" binding:"required"`
	Priority      string    `json:"priority"`
	DueDate       time.Time `json:"dueDate"`
	AppliedBy     string    `json:"appliedBy" binding:"required"`
	AssesseeID    uuid.UUID `json:"assesseeId" binding:"required"`
	IsDraft       bool      `json:"isDraft"`
	ImportantNote string    `json:"importantNote" binding:"max=500,omitempty"`
	// PropertyOwners []PropertyOwnerRequest `json:"propertyOwners" binding:"required,min=1,dive"`
}

// PropertyOwnerRequest holds property owner details in an application request
type PropertyOwnerRequest struct {
	AdhaarNo               uint64  `json:"adhaarNo" binding:"required"`
	Name                   string  `json:"name" binding:"required"`
	ContactNo              string  `json:"contactNo" binding:"required"`
	Email                  string  `json:"email" binding:"omitempty,email"`
	Gender                 string  `json:"gender" binding:"required,oneof=MALE FEMALE OTHER"`
	Guardian               string  `json:"guardian"`
	GuardianType           string  `json:"guardianType" binding:"omitempty,oneof=FATHER MOTHER SPOUSE GUARDIAN OTHER"`
	RelationshipToProperty string  `json:"relationshipToProperty" binding:"omitempty,oneof=OWNER CO_OWNER JOINT_OWNER LEGAL_HEIR POWER_OF_ATTORNEY TENANT FAMILY_MEMBER OTHER"`
	OwnershipShare         float64 `json:"ownershipShare" binding:"required,min=0.01,max=100"`
	IsPrimaryOwner         bool    `json:"isPrimaryOwner"`
}

// UpdateApplicationRequest is used to update an existing application
type UpdateApplicationRequest struct {
	Priority       string                 `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH"`
	DueDate        *time.Time             `json:"dueDate"`
	AssignedAgent  *uuid.UUID             `json:"assignedAgent"`
	AppliedBy      string                 `json:"appliedBy"`
	IsDraft        *bool                  `json:"isDraft"`
	PropertyOwners []PropertyOwnerRequest `json:"propertyOwners" binding:"omitempty,dive"`
	ImportantNote  *string                `json:"importantNote" binding:"omitempty,max=500"`
}

// UpdateStatusRequest is used to update the status of an application
type UpdateStatusRequest struct {
	Status   string  `json:"status" binding:"required"`
	Comments *string `json:"comments"`
}

// ActionRequest is used for agent assignment and action requests
type ActionRequest struct {
	Action   string    `json:"action"`
	AgentID  uuid.UUID `json:"agentId"`
	Comments string    `json:"comments"`
	Verified bool      `json:"verified"`
	Approved bool      `json:"approved"`
}

// ApplicationResponse represents the application response
type ApplicationResponse struct {
	ID                 uuid.UUID  `json:"id"`
	ApplicationNo      string     `json:"applicationNo"`
	PropertyID         uuid.UUID  `json:"propertyId"`
	Priority           string     `json:"priority"`
	DueDate            time.Time  `json:"dueDate"`
	AssignedAgent      *uuid.UUID `json:"assignedAgent,omitempty"`
	Status             string     `json:"status"`
	WorkflowInstanceID string     `json:"workflowInstanceId,omitempty"`
	AppliedBy          string     `json:"appliedBy"`
	AssesseeID         *uuid.UUID `json:"assesseeId,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	ImportantNote      string     `json:"importantNote"`
}

// ApplicationSearchCriteria holds filters for searching applications
type ApplicationSearchCriteria struct {
	Status          string     `form:"status"`
	Priority        string     `form:"priority"`
	PropertyID      string     `form:"propertyId"`
	AssignedAgent   string     `form:"assignedAgent"`
	AppliedBy       string     `form:"appliedBy"`
	ApplicationNo   string     `form:"applicationNo"`
	CreatedDateFrom *time.Time `form:"createdDateFrom"`
	CreatedDateTo   *time.Time `form:"createdDateTo"`
	DueDateFrom     *time.Time `form:"dueDateFrom"`
	DueDateTo       *time.Time `form:"dueDateTo"`
	IsDraft         *bool      `form:"isDraft"`
	ZoneNo          string     `form:"zoneNo"`
	WardNo          []string   `form:"wardNo"`
	AssesseeID      string     `form:"assesseeId"`
	SortBy          string     `form:"sortBy"`
	SortField       string     `form:"sortField"` // created_at or due_date
	IsCountOnly     bool       `form:"isCountOnly"`
	PropertyNo      string     `form:"propertyNo"`
}

// ============================================================================
// PROPERTY OWNER DTOs
// ============================================================================

// CreatePropertyOwnerRequest is used to add a new property owner
type CreatePropertyOwnerRequest struct {
	PropertyID             uuid.UUID `json:"propertyId" binding:"required"`
	AdhaarNo               uint64    `json:"adhaarNo" binding:"required"`
	Name                   string    `json:"name" binding:"required"`
	ContactNo              string    `json:"contactNo" binding:"required"`
	Email                  string    `json:"email" binding:"omitempty,email"`
	Gender                 string    `json:"gender" binding:"required,oneof=MALE FEMALE OTHER"`
	Guardian               string    `json:"guardian"`
	GuardianType           string    `json:"guardianType" binding:"omitempty,oneof=FATHER MOTHER SPOUSE OTHER"`
	RelationshipToProperty string    `json:"relationshipToProperty" binding:"omitempty,oneof=OWNER CO_OWNER JOINT_OWNER LEGAL_HEIR POWER_OF_ATTORNEY OTHER"`
	OwnershipShare         float64   `json:"ownershipShare" binding:"required,min=0.01,max=100"`
	IsPrimaryOwner         bool      `json:"isPrimaryOwner"`
}

// PropertyOwnerResponse is returned as the response for property owner APIs
type PropertyOwnerResponse struct {
	ID                     uuid.UUID `json:"id"`
	ApplicationID          string    `json:"applicationId"`
	PropertyID             uuid.UUID `json:"propertyId"`
	AdhaarNo               uint64    `json:"adhaarNo"`
	Name                   string    `json:"name"`
	ContactNo              string    `json:"contactNo"`
	Email                  string    `json:"email,omitempty"`
	Gender                 string    `json:"gender"`
	Guardian               string    `json:"guardian,omitempty"`
	GuardianType           string    `json:"guardianType,omitempty"`
	RelationshipToProperty string    `json:"relationshipToProperty,omitempty"`
	OwnershipShare         float64   `json:"ownershipShare"`
	IsPrimaryOwner         bool      `json:"isPrimaryOwner"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

// UpdatePropertyOwnerRequest is used to update property owner details
type UpdatePropertyOwnerRequest struct {
	Name                   string  `json:"name"`
	ContactNo              string  `json:"contactNo"`
	Email                  string  `json:"email" binding:"omitempty,email"`
	Gender                 string  `json:"gender" binding:"omitempty,oneof=MALE FEMALE OTHER"`
	Guardian               string  `json:"guardian"`
	GuardianType           string  `json:"guardianType" binding:"omitempty,oneof=FATHER MOTHER SPOUSE GUARDIAN OTHER"`
	RelationshipToProperty string  `json:"relationshipToProperty" binding:"omitempty,oneof=OWNER CO_OWNER JOINT_OWNER LEGAL_HEIR POWER_OF_ATTORNEY TENANT FAMILY_MEMBER OTHER"`
	OwnershipShare         float64 `json:"ownershipShare" binding:"omitempty,min=0.01,max=100"`
	IsPrimaryOwner         bool    `json:"isPrimaryOwner"`
}

//IGRS DTOs

type CreateIGRSRequest struct {
	PropertyID         uuid.UUID `json:"propertyId" binding:"required"`
	Habitation         string    `json:"habitation" binding:"required"`
	IGRSWard           string    `json:"igrsWard" binding:"omitempty"`
	IGRSLocality       string    `json:"igrsLocality" binding:"omitempty"`
	IGRSBlock          string    `json:"igrsBlock" binding:"omitempty"`
	DoorNoFrom         string    `json:"doorNoFrom" binding:"omitempty"`
	DoorNoTo           string    `json:"doorNoTo" binding:"omitempty"`
	IGRSClassification string    `json:"igrsClassification" binding:"omitempty"`
	BuiltUpAreaPct     *float64  `json:"builtUpAreaPct" binding:"omitempty,gte=0,lte=100"`
	FrontSetback       *float64  `json:"frontSetback" binding:"omitempty,gte=0"`
	RearSetback        *float64  `json:"rearSetback" binding:"omitempty,gte=0"`
	SideSetback        *float64  `json:"sideSetback" binding:"omitempty,gte=0"`
	TotalPlinthArea    *float64  `json:"totalPlinthArea" binding:"omitempty,gte=0"`
}

type UpdateIGRSRequest = CreateIGRSRequest

// CreateApplicationLogRequest represents request for creating application log
type CreateApplicationLogRequest struct {
	Action        string                 `json:"action"`
	PerformedBy   string                 `json:"performedBy" binding:"required"`
	Comments      string                 `json:"comments"`
	Metadata      map[string]interface{} `json:"metadata"`
	FileStoreID   *uuid.UUID             `json:"fileStoreId"`
	Actor         string                 `json:"actor"`
	ApplicationID uuid.UUID              `json:"applicationId" binding:"required"`
}

// UpdateApplicationLogRequest represents request for updating application log
type UpdateApplicationLogRequest struct {
	Action      string                 `json:"action,omitempty" binding:"omitempty,oneof=CREATE UPDATE FILE_UPLOAD STATUS_CHANGE"`
	PerformedBy string                 `json:"performedBy,omitempty"`
	Comments    string                 `json:"comments,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	FileStoreID *uuid.UUID             `json:"fileStoreId,omitempty"`
}
