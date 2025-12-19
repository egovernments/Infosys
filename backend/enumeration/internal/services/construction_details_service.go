package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

// Ensure constructionDetailsService implements ConstructionDetailsService interface at compile time.
var _ ConstructionDetailsService = (*constructionDetailsService)(nil)

// constructionDetailsService provides business logic for managing ConstructionDetails entities using a repository.
type constructionDetailsService struct {
	constructionDetailsRepo repositories.ConstructionDetailsRepository
}

// NewConstructionDetailsService creates a new instance of constructionDetailsService.
// constructionDetailsRepo: repository for ConstructionDetails data access.
// Returns: ConstructionDetailsService implementation.
func NewConstructionDetailsService(constructionDetailsRepo repositories.ConstructionDetailsRepository) ConstructionDetailsService {
	return &constructionDetailsService{
		constructionDetailsRepo: constructionDetailsRepo,
	}
}

// CreateConstructionDetails validates and creates a new ConstructionDetails record.
// ctx: context for the operation.
// constructionDetails: pointer to ConstructionDetails model to be created.
// Returns: error if validation or creation fails.
func (s *constructionDetailsService) CreateConstructionDetails(ctx context.Context, constructionDetails *models.ConstructionDetails) error {
	if constructionDetails == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if err := s.validateConstructionDetails(constructionDetails); err != nil {
		return err // Service layer validation error - return as is
	}
	return s.constructionDetailsRepo.Create(ctx, constructionDetails)
}

// GetConstructionDetailsByID retrieves a ConstructionDetails record by its unique ID.
// ctx: context for the operation.
// id: UUID of the construction details.
// Returns: pointer to ConstructionDetails and error if not found or on failure.
func (s *constructionDetailsService) GetConstructionDetailsByID(ctx context.Context, id uuid.UUID) (*models.ConstructionDetails, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid construction details ID: cannot be nil")
	}
	// Repository layer error - propagate as is
	return s.constructionDetailsRepo.GetByID(ctx, id)
}

// UpdateConstructionDetails validates and updates an existing ConstructionDetails record.
// ctx: context for the operation.
// constructionDetails: pointer to ConstructionDetails model with updated data.
// Returns: error if validation or update fails.
func (s *constructionDetailsService) UpdateConstructionDetails(ctx context.Context, constructionDetails *models.ConstructionDetails) error {
	if constructionDetails == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if constructionDetails.ID == uuid.Nil {
		return fmt.Errorf("invalid construction details ID: cannot be nil")
	}

	if err := s.validateConstructionDetails(constructionDetails); err != nil {
		return err // Service layer validation error - return as is
	}

	// Repository layer error - propagate as is
	return s.constructionDetailsRepo.Update(ctx, constructionDetails)
}

// DeleteConstructionDetails removes a ConstructionDetails record by its unique ID.
// ctx: context for the operation.
// id: UUID of the construction details to delete.
// Returns: error if deletion fails or record not found.
func (s *constructionDetailsService) DeleteConstructionDetails(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid construction details ID: cannot be nil")
	}
	// Repository layer error - propagate as is
	return s.constructionDetailsRepo.Delete(ctx, id)
}

// GetAllConstructionDetails retrieves all ConstructionDetails records, optionally filtered by propertyID, with pagination.
// ctx: context for the operation.
// page: page number (zero-based), size: number of records per page, propertyID: optional filter.
// Returns: slice of ConstructionDetails pointers, total count, and error if any.
func (s *constructionDetailsService) GetAllConstructionDetails(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.ConstructionDetails, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}
	return s.constructionDetailsRepo.GetAll(ctx, page, size, propertyID)
}

// GetConstructionDetailsByPropertyID retrieves ConstructionDetails records by their associated property ID.
// ctx: context for the operation.
// propertyID: UUID of the property.
// Returns: slice of ConstructionDetails pointers and error if not found or on failure.
func (s *constructionDetailsService) GetConstructionDetailsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.ConstructionDetails, error) {
	if propertyID == uuid.Nil {
		return nil, fmt.Errorf("invalid property ID: cannot be nil")
	}
	return s.constructionDetailsRepo.GetByPropertyID(ctx, propertyID)
}

// validateConstructionDetails checks required fields for ConstructionDetails and returns an error if validation fails.
// constructionDetails: pointer to ConstructionDetails model to validate.
// Returns: error if validation fails, nil otherwise.
func (s *constructionDetailsService) validateConstructionDetails(constructionDetails *models.ConstructionDetails) error {
	if constructionDetails.PropertyID == uuid.Nil {
		return fmt.Errorf("property ID is required")
	}

	// Optional validation - you can customize based on your business rules
	if constructionDetails.FloorType == "" {
		return fmt.Errorf("floor type is required")
	}

	if constructionDetails.WallType == "" {
		return fmt.Errorf("wall type is required")
	}

	if constructionDetails.RoofType == "" {
		return fmt.Errorf("roof type is required")
	}

	return nil
}
