package validators

import (
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Unit tests for CoordinatesValidator covering validation of coordinates, geographic range, batch operations, and error formatting

func TestCoordinatesValidator_ValidateCoordinates(t *testing.T) {
	// Test various cases for validating a single coordinates object
	validator := NewCoordinatesValidator()
	gisDataID := uuid.New()

	// Each test case checks a different validation scenario
	tests := []struct {
		name        string
		coordinates *models.Coordinates
		wantErr     bool
		errField    string
	}{
		{
			name: "valid coordinates",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name:        "nil coordinates",
			coordinates: nil,
			wantErr:     true,
			errField:    "coordinates",
		},
		{
			name: "latitude too low",
			coordinates: &models.Coordinates{
				Latitude:  -91.0,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "latitude",
		},
		{
			name: "latitude too high",
			coordinates: &models.Coordinates{
				Latitude:  91.0,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "latitude",
		},
		{
			name: "longitude too low",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: -181.0,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "longitude",
		},
		{
			name: "longitude too high",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 181.0,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "longitude",
		},
		{
			name: "both zero coordinates",
			coordinates: &models.Coordinates{
				Latitude:  0.0,
				Longitude: 0.0,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "coordinates",
		},
		{
			name: "missing gisDataId",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: uuid.Nil,
			},
			wantErr:  true,
			errField: "gisDataId",
		},
		{
			name: "edge case: min valid latitude",
			coordinates: &models.Coordinates{
				Latitude:  constants.MinLatitude,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name: "edge case: max valid latitude",
			coordinates: &models.Coordinates{
				Latitude:  constants.MaxLatitude,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name: "edge case: min valid longitude",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: constants.MinLongitude,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name: "edge case: max valid longitude",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: constants.MaxLongitude,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
	}

	// Run all test cases for ValidateCoordinates
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCoordinates(tt.coordinates)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errField != "" {
					valErr, ok := err.(*ValidationError)
					assert.True(t, ok, "expected ValidationError")
					if ok {
						assert.Equal(t, tt.errField, valErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validation of coordinates within Karnataka's geographic range
func TestCoordinatesValidator_ValidateGeographicRange(t *testing.T) {
	validator := NewCoordinatesValidator()
	gisDataID := uuid.New()

	// Each test case checks a different geographic range scenario
	tests := []struct {
		name        string
		coordinates *models.Coordinates
		wantErr     bool
		errField    string
	}{
		{
			name: "within Karnataka bounds - Bangalore",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name: "within Karnataka bounds - Mangalore",
			coordinates: &models.Coordinates{
				Latitude:  12.9141,
				Longitude: 74.8560,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name:        "nil coordinates",
			coordinates: nil,
			wantErr:     true,
			errField:    "coordinates",
		},
		{
			name: "latitude below Karnataka minimum",
			coordinates: &models.Coordinates{
				Latitude:  10.0,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "latitude",
		},
		{
			name: "latitude above Karnataka maximum",
			coordinates: &models.Coordinates{
				Latitude:  19.0,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "latitude",
		},
		{
			name: "longitude below Karnataka minimum",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 73.0,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "longitude",
		},
		{
			name: "longitude above Karnataka maximum",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 79.0,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "longitude",
		},
		{
			name: "edge case: min Karnataka latitude",
			coordinates: &models.Coordinates{
				Latitude:  constants.KarnatakaMinLatitude,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name: "edge case: max Karnataka latitude",
			coordinates: &models.Coordinates{
				Latitude:  constants.KarnatakaMaxLatitude,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
	}

	// Run all test cases for ValidateGeographicRange
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateGeographicRange(tt.coordinates)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errField != "" {
					valErr, ok := err.(*ValidationError)
					assert.True(t, ok, "expected ValidationError")
					if ok {
						assert.Equal(t, tt.errField, valErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test batch validation of coordinates
func TestCoordinatesValidator_ValidateBatch(t *testing.T) {
	validator := NewCoordinatesValidator()
	gisDataID := uuid.New()

	// Each test case checks a different batch validation scenario
	tests := []struct {
		name        string
		coordinates []*models.Coordinates
		wantErr     bool
	}{
		{
			name: "valid batch",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				{Latitude: 13.0827, Longitude: 77.5828, GISDataID: gisDataID},
			},
			wantErr: false,
		},
		{
			name: "valid batch - single item",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
			},
			wantErr: false,
		},
		{
			name:        "nil batch",
			coordinates: nil,
			wantErr:     true,
		},
		{
			name:        "empty batch",
			coordinates: []*models.Coordinates{},
			wantErr:     true,
		},
		{
			name: "batch with nil coordinate",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				nil,
			},
			wantErr: true,
		},
		{
			name: "batch with invalid latitude",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				{Latitude: 91.0, Longitude: 77.5946, GISDataID: gisDataID},
			},
			wantErr: true,
		},
		{
			name: "batch with invalid longitude",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				{Latitude: 12.9716, Longitude: 181.0, GISDataID: gisDataID},
			},
			wantErr: true,
		},
		{
			name: "batch with missing GISDataID",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				{Latitude: 13.0827, Longitude: 77.5828, GISDataID: uuid.Nil},
			},
			wantErr: true,
		},
		{
			name: "batch with multiple errors",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				nil,
				{Latitude: 91.0, Longitude: 181.0, GISDataID: uuid.Nil},
			},
			wantErr: true,
		},
	}

	// Run all test cases for ValidateBatch
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBatch(tt.coordinates)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test batch validation of coordinates with geographic range
func TestCoordinatesValidator_ValidateBatchWithGeographicRange(t *testing.T) {
	validator := NewCoordinatesValidator()
	gisDataID := uuid.New()

	// Each test case checks a different batch validation scenario for geographic range
	tests := []struct {
		name        string
		coordinates []*models.Coordinates
		wantErr     bool
	}{
		{
			name: "valid batch within Karnataka",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				{Latitude: 13.0827, Longitude: 77.5828, GISDataID: gisDataID},
			},
			wantErr: false,
		},
		{
			name: "batch with coordinates outside Karnataka",
			coordinates: []*models.Coordinates{
				{Latitude: 12.9716, Longitude: 77.5946, GISDataID: gisDataID},
				{Latitude: 10.0, Longitude: 77.5946, GISDataID: gisDataID},
			},
			wantErr: true,
		},
		{
			name:        "nil batch",
			coordinates: nil,
			wantErr:     true,
		},
		{
			name:        "empty batch",
			coordinates: []*models.Coordinates{},
			wantErr:     true,
		},
	}

	// Run all test cases for ValidateBatchWithGeographicRange
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBatchWithGeographicRange(tt.coordinates)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validation of update requests for coordinates
func TestCoordinatesValidator_ValidateUpdateRequest(t *testing.T) {
	validator := NewCoordinatesValidator()
	gisDataID := uuid.New()
	coordID := uuid.New()

	// Each test case checks a different update validation scenario
	tests := []struct {
		name        string
		coordinates *models.Coordinates
		wantErr     bool
		errField    string
	}{
		{
			name: "valid update",
			coordinates: &models.Coordinates{
				ID:        coordID,
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name:        "nil coordinates",
			coordinates: nil,
			wantErr:     true,
			errField:    "coordinates",
		},
		{
			name: "missing ID",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "id",
		},
		{
			name: "invalid latitude with valid ID",
			coordinates: &models.Coordinates{
				ID:        coordID,
				Latitude:  91.0,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "latitude",
		},
		{
			name: "invalid longitude with valid ID",
			coordinates: &models.Coordinates{
				ID:        coordID,
				Latitude:  12.9716,
				Longitude: 181.0,
				GISDataID: gisDataID,
			},
			wantErr:  true,
			errField: "longitude",
		},
		{
			name: "missing GISDataID with valid ID",
			coordinates: &models.Coordinates{
				ID:        coordID,
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: uuid.Nil,
			},
			wantErr:  true,
			errField: "gisDataId",
		},
	}

	// Run all test cases for ValidateUpdateRequest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdateRequest(tt.coordinates)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errField != "" {
					valErr, ok := err.(*ValidationError)
					assert.True(t, ok, "expected ValidationError")
					if ok {
						assert.Equal(t, tt.errField, valErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test validation of create requests for coordinates
func TestCoordinatesValidator_ValidateRequest(t *testing.T) {
	validator := NewCoordinatesValidator()
	gisDataID := uuid.New()

	// Each test case checks a different create request scenario
	tests := []struct {
		name        string
		coordinates *models.Coordinates
		wantErr     bool
	}{
		{
			name: "valid request",
			coordinates: &models.Coordinates{
				Latitude:  12.9716,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: false,
		},
		{
			name:        "nil request",
			coordinates: nil,
			wantErr:     true,
		},
		{
			name: "invalid request",
			coordinates: &models.Coordinates{
				Latitude:  91.0,
				Longitude: 77.5946,
				GISDataID: gisDataID,
			},
			wantErr: true,
		},
	}

	// Run all test cases for ValidateRequest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateRequest(tt.coordinates)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test error formatting for ValidationError
func TestValidationError_Error(t *testing.T) {
	// Each test case checks error message formatting
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "error with value",
			err: &ValidationError{
				Field:   "latitude",
				Message: "must be between -90 and 90",
				Value:   91.0,
			},
			expected: "validation failed for field 'latitude': must be between -90 and 90 (received: 91)",
		},
		{
			name: "error without value",
			err: &ValidationError{
				Field:   "coordinates",
				Message: "coordinates object cannot be nil",
				Value:   nil,
			},
			expected: "validation failed for field 'coordinates': coordinates object cannot be nil",
		},
	}

	// Run all test cases for ValidationError error formatting
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

// Test error formatting for BatchValidationError
func TestBatchValidationError_Error(t *testing.T) {
	err1 := NewValidationError("latitude", "invalid", 91.0)
	err2 := NewValidationError("longitude", "invalid", 181.0)

	batchErr := NewBatchValidationError([]error{err1, err2})
	expected := "batch validation failed with 2 errors"

	assert.Equal(t, expected, batchErr.Error())
	assert.Len(t, batchErr.Errors, 2)
}
