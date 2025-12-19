// Package constants holds error message constants for the service
package constants

// Error message constants for various failure scenarios
const (
    // Application Errors
    ErrApplicationNotFound           = "application not found"
    ErrApplicationAlreadyExists      = "application already exists"
    ErrApplicationInvalidStatus      = "invalid application status"
    ErrApplicationInvalidPriority    = "invalid application priority"
    ErrApplicationCreationFailed     = "failed to create application"
    ErrApplicationUpdateFailed       = "failed to update application" 
    ErrApplicationDeletionFailed     = "failed to delete application"
    
    // Workflow Errors
    ErrWorkflowProcessNotFound       = "workflow process not found - please ensure workflow is pre-created"
    ErrWorkflowInstanceCreationFailed = "failed to create workflow instance"
    
    // Validation Errors
    ErrInvalidApplicationID          = "invalid application ID format"
    ErrInvalidPropertyID            = "invalid property ID"
    ErrInvalidTenantID              = "invalid tenant ID"
    ErrInvalidCitizenID             = "invalid citizen ID"
    ErrMissingRequiredFields        = "missing required fields"
    ErrExceededMaxApplications      = "exceeded maximum applications per citizen"
    
    // Property Errors
    ErrPropertyNotFound             = "property not found"
    ErrPropertyValidationFailed     = "property validation failed"
    
    // Coordinates Errors
    ErrCoordinatesNotFound          = "coordinates not found"
    ErrCoordinatesValidationFailed  = "coordinates validation failed"
    ErrInvalidLatitude              = "invalid latitude value"
    ErrInvalidLongitude             = "invalid longitude value"
    ErrInvalidGISDataID             = "invalid GIS data ID"
    ErrCoordinatesOutOfBounds       = "coordinates out of geographic bounds"
    ErrBatchValidationFailed        = "batch validation failed"
    
    // General Errors
    ErrInternalServer              = "internal server error"
    ErrUnauthorized               = "unauthorized access"
    ErrForbidden                  = "forbidden access"
    ErrBadRequest                 = "bad request"
)
