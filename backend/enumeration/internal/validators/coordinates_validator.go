package validators

import (
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
)

// CoordinatesValidator provides validation logic for coordinates
type CoordinatesValidator struct{}

// NewCoordinatesValidator returns a new CoordinatesValidator
func NewCoordinatesValidator() *CoordinatesValidator {
	return &CoordinatesValidator{}
}

// ValidateCoordinates checks if a single coordinates object is valid
func (v *CoordinatesValidator) ValidateCoordinates(coordinates *models.Coordinates) error {
	// Nil pointer check
	if coordinates == nil {
		return NewValidationError("coordinates", "coordinates object cannot be nil", nil)
	}

	// Validate latitude range
	if coordinates.Latitude < constants.MinLatitude || coordinates.Latitude > constants.MaxLatitude {
		return NewValidationError(
			"latitude",
			fmt.Sprintf("must be between %.1f and %.1f", constants.MinLatitude, constants.MaxLatitude),
			coordinates.Latitude,
		)
	}

	// Validate longitude range
	if coordinates.Longitude < constants.MinLongitude || coordinates.Longitude > constants.MaxLongitude {
		return NewValidationError(
			"longitude",
			fmt.Sprintf("must be between %.1f and %.1f", constants.MinLongitude, constants.MaxLongitude),
			coordinates.Longitude,
		)
	}

	// Check for exact zero values (potential uninitialized coordinates)
	if coordinates.Latitude == 0.0 && coordinates.Longitude == 0.0 {
		return NewValidationError(
			"coordinates",
			"latitude and longitude cannot both be zero (0.0, 0.0)",
			fmt.Sprintf("(%.6f, %.6f)", coordinates.Latitude, coordinates.Longitude),
		)
	}

	// Check if GISDataID is provided
	if coordinates.GISDataID == uuid.Nil {
		return NewValidationError("gisDataId", "gisDataId is required", coordinates.GISDataID)
	}

	return nil
}

// ValidateGeographicRange checks if coordinates are within India's bounds
func (v *CoordinatesValidator) ValidateGeographicRange(coordinates *models.Coordinates) error {
	if coordinates == nil {
		return NewValidationError("coordinates", "coordinates object cannot be nil", nil)
	}

	// Check India bounds
	if coordinates.Latitude < constants.IndiaMinLatitude || coordinates.Latitude > constants.IndiaMaxLatitude {
		return NewValidationError(
			"latitude",
			fmt.Sprintf("must be within India bounds (%.1f to %.1f)", constants.IndiaMinLatitude, constants.IndiaMaxLatitude),
			coordinates.Latitude,
		)
	}

	if coordinates.Longitude < constants.IndiaMinLongitude || coordinates.Longitude > constants.IndiaMaxLongitude {
		return NewValidationError(
			"longitude",
			fmt.Sprintf("must be within India bounds (%.1f to %.1f)", constants.IndiaMinLongitude, constants.IndiaMaxLongitude),
			coordinates.Longitude,
		)
	}

	return nil
}

// ValidateKarnatakaRange checks if coordinates are within Karnataka's bounds
func (v *CoordinatesValidator) ValidateKarnatakaRange(coordinates *models.Coordinates) error {
	if coordinates == nil {
		return NewValidationError("coordinates", "coordinates object cannot be nil", nil)
	}

	// Check Karnataka bounds
	if coordinates.Latitude < constants.KarnatakaMinLatitude || coordinates.Latitude > constants.KarnatakaMaxLatitude {
		return NewValidationError(
			"latitude",
			fmt.Sprintf("must be within Karnataka bounds (%.1f to %.1f)", constants.KarnatakaMinLatitude, constants.KarnatakaMaxLatitude),
			coordinates.Latitude,
		)
	}

	if coordinates.Longitude < constants.KarnatakaMinLongitude || coordinates.Longitude > constants.KarnatakaMaxLongitude {
		return NewValidationError(
			"longitude",
			fmt.Sprintf("must be within Karnataka bounds (%.1f to %.1f)", constants.KarnatakaMinLongitude, constants.KarnatakaMaxLongitude),
			coordinates.Longitude,
		)
	}

	return nil
}

// ValidateRequest checks if a create request for coordinates is valid
func (v *CoordinatesValidator) ValidateRequest(coordinates *models.Coordinates) error {
	return v.ValidateCoordinates(coordinates)
}

// ValidateUpdateRequest checks if an update request for coordinates is valid
func (v *CoordinatesValidator) ValidateUpdateRequest(coordinates *models.Coordinates) error {
	if coordinates == nil {
		return NewValidationError("coordinates", "coordinates object cannot be nil", nil)
	}

	// Check ID for update
	if coordinates.ID == uuid.Nil {
		return NewValidationError("id", "id is required for update operation", coordinates.ID)
	}

	// Validate coordinate values
	return v.ValidateCoordinates(coordinates)
}

// ValidateBatch checks if a batch of coordinates is valid
func (v *CoordinatesValidator) ValidateBatch(coordinates []*models.Coordinates) error {
	if coordinates == nil {
		return NewValidationError("coordinates", "coordinates array cannot be nil", nil)
	}

	if len(coordinates) == 0 {
		return NewValidationError("coordinates", "coordinates array cannot be empty", len(coordinates))
	}

	var errors []error
	for i, coord := range coordinates {
		if coord == nil {
			errors = append(errors, NewValidationError(
				fmt.Sprintf("coordinates[%d]", i),
				"coordinate cannot be nil",
				nil,
			))
			continue
		}

		if err := v.ValidateCoordinates(coord); err != nil {
			errors = append(errors, fmt.Errorf("coordinates[%d]: %w", i, err))
		}
	}

	if len(errors) > 0 {
		return NewBatchValidationError(errors)
	}

	return nil
}

// ValidateBatchWithGeographicRange checks if a batch of coordinates is valid and within India's bounds
func (v *CoordinatesValidator) ValidateBatchWithGeographicRange(coordinates []*models.Coordinates) error {
	if err := v.ValidateBatch(coordinates); err != nil {
		return err
	}

	var errors []error
	for i, coord := range coordinates {
		if coord == nil {
			continue
		}

		if err := v.ValidateGeographicRange(coord); err != nil {
			errors = append(errors, fmt.Errorf("coordinates[%d]: %w", i, err))
		}
	}

	if len(errors) > 0 {
		return NewBatchValidationError(errors)
	}

	return nil
}
