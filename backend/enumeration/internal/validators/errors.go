package validators

import "fmt"

// ValidationError describes a single validation error for a field
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("validation failed for field '%s': %s (received: %v)", e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NewValidationError returns a new ValidationError
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// BatchValidationError groups multiple validation errors together
type BatchValidationError struct {
	Errors []error
}

func (e *BatchValidationError) Error() string {
	return fmt.Sprintf("batch validation failed with %d errors", len(e.Errors))
}

// NewBatchValidationError returns a new BatchValidationError
func NewBatchValidationError(errors []error) *BatchValidationError {
	return &BatchValidationError{Errors: errors}
}
