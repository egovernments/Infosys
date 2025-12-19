package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"enumeration/internal/validators"
	"fmt"

	"github.com/google/uuid"
)

// Compile-time check for interface implementation
var _ CoordinatesService = (*coordinatesService)(nil)

// coordinatesService handles all business logic for coordinates
type coordinatesService struct {
	repo      repositories.CoordinatesRepository
	validator *validators.CoordinatesValidator
}

// NewCoordinatesService returns a new coordinatesService
func NewCoordinatesService(repo repositories.CoordinatesRepository) CoordinatesService {
	return &coordinatesService{
		repo:      repo,
		validator: validators.NewCoordinatesValidator(),
	}
}

// FindAll gets coordinates with pagination and optional GISDataID filter
func (s *coordinatesService) FindAll(ctx context.Context, page, size int, gisDataID *uuid.UUID) ([]models.Coordinates, int64, error) {
	// Nil context check
	return s.repo.FindAll(ctx, page, size, gisDataID)
}

// Create adds a new coordinates record after validation
func (s *coordinatesService) Create(ctx context.Context, coordinates *models.Coordinates) error {
	// Nil checks
	if coordinates == nil {
		return validators.NewValidationError("coordinates", "coordinates object cannot be nil", nil)
	}

	// Validate request
	if err := s.validator.ValidateRequest(coordinates); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create in repository
	if err := s.repo.Create(ctx, coordinates); err != nil {
		return fmt.Errorf("failed to create coordinates for GISDataID %s: %w", coordinates.GISDataID, err)
	}

	return nil
}

// GetByID fetches a coordinates record by its ID
func (s *coordinatesService) GetByID(ctx context.Context, id uuid.UUID) (*models.Coordinates, error) {
	if id == uuid.Nil {
		return nil, validators.NewValidationError("id", "id cannot be nil", id)
	}

	// Fetch from repository
	coordinates, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates with ID %s: %w", id, err)
	}

	return coordinates, nil
}

// Update modifies an existing coordinates record after validation
func (s *coordinatesService) Update(ctx context.Context, coordinates *models.Coordinates) error {
	// Nil checks
	if coordinates == nil {
		return validators.NewValidationError("coordinates", "coordinates object cannot be nil", nil)
	}

	// Validate update request
	if err := s.validator.ValidateUpdateRequest(coordinates); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update in repository
	if err := s.repo.Update(ctx, coordinates); err != nil {
		return fmt.Errorf("failed to update coordinates with ID %s: %w", coordinates.ID, err)
	}

	return nil
}

// Delete removes a coordinates record by its ID
func (s *coordinatesService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return validators.NewValidationError("id", "id cannot be nil", id)
	}

	// Delete from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete coordinates with ID %s: %w", id, err)
	}

	return nil
}

// GetAll returns all coordinates with pagination, as pointers
func (s *coordinatesService) GetAll(ctx context.Context, page, size int, gisDataID *uuid.UUID) ([]*models.Coordinates, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 {
		size = constants.DefaultSize
	}

	// Fetch from repository
	coordinates, total, err := s.repo.FindAll(ctx, page, size, gisDataID)
	if err != nil {
		gisIDStr := "all"
		if gisDataID != nil {
			gisIDStr = gisDataID.String()
		}
		return nil, 0, fmt.Errorf("failed to fetch coordinates for GISDataID %s (page: %d, size: %d): %w", gisIDStr, page, size, err)
	}

	// Convert []models.Coordinates to []*models.Coordinates
	result := make([]*models.Coordinates, len(coordinates))
	for i := range coordinates {
		result[i] = &coordinates[i]
	}

	return result, total, nil
}

// CreateBatch adds multiple coordinates records after validation
func (s *coordinatesService) CreateBatch(ctx context.Context, coords []*models.Coordinates) error {
	// Nil checks
	if coords == nil {
		return validators.NewValidationError("coordinates", "coordinates array cannot be nil", nil)
	}
	if len(coords) == 0 {
		return validators.NewValidationError("coordinates", "coordinates array cannot be empty", len(coords))
	}

	// Validate batch
	if err := s.validator.ValidateBatch(coords); err != nil {
		return fmt.Errorf("batch validation failed: %w", err)
	}

	// Create batch in repository
	if err := s.repo.CreateBatch(ctx, coords); err != nil {
		return fmt.Errorf("failed to create batch of %d coordinates: %w", len(coords), err)
	}

	return nil
}

// ReplaceByGISDataID replaces all coordinates for a GISDataID with the provided batch
func (s *coordinatesService) ReplaceByGISDataID(ctx context.Context, gisDataID uuid.UUID, coords []*models.Coordinates) error {
	if gisDataID == uuid.Nil {
		return validators.NewValidationError("gisDataId", "gisDataId cannot be nil", gisDataID)
	}

	// Allow empty slice (delete all coordinates for gisDataID)
	if len(coords) == 0 {
		if err := s.repo.ReplaceByGISDataID(ctx, gisDataID, coords); err != nil {
			return fmt.Errorf("failed to clear coordinates for GISDataID %s: %w", gisDataID, err)
		}
		return nil
	}

	// Normalize and validate each coordinate
	for i, c := range coords {
		if c == nil {
			return validators.NewValidationError(
				fmt.Sprintf("coordinates[%d]", i),
				"coordinate cannot be nil",
				nil,
			)
		}

		// Ensure the coordinate points to the requested GISDataID
		c.GISDataID = gisDataID

		// Assign ID if missing
		if c.ID == uuid.Nil {
			c.ID = uuid.New()
		}
	}

	// Validate batch
	if err := s.validator.ValidateBatch(coords); err != nil {
		return fmt.Errorf("batch validation failed for GISDataID %s: %w", gisDataID, err)
	}

	// Replace in repository (transactional)
	if err := s.repo.ReplaceByGISDataID(ctx, gisDataID, coords); err != nil {
		return fmt.Errorf("failed to replace %d coordinates for GISDataID %s: %w", len(coords), gisDataID, err)
	}

	return nil
}
