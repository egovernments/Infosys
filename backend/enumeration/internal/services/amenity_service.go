// Package services provides business logic and service layer implementations
// for the property tax enumeration system.
package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"enumeration/pkg/logger"
	"fmt"
)

// Compile-time check to ensure amenityService implements AmenityService interface.
var _ AmenityService = (*amenityService)(nil)

// amenityService implements the AmenityService interface and provides
// methods for managing amenity records.
type amenityService struct {
	repo repositories.AmenityRepository // Repository for amenity data access
}

// NewAmenityService creates a new instance of AmenityService with the provided repository.
func NewAmenityService(repo repositories.AmenityRepository) AmenityService {
	return &amenityService{repo}
}

// GetAll returns all amenities with pagination.
// Service enforces sane defaults and max page size.
func (s *amenityService) GetAll(ctx context.Context, page, size int) ([]models.Amenities, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 {
		size = constants.DefaultSize
	}
	if size > 100 {
		size = constants.MaxSize
	}
	// Reuse repository filter method with no filters
	amenities, total, err := s.repo.GetAllWithFilters(ctx, page, size, "", "")
	if err != nil {
		return nil, 0, err
	}
	logger.Info("Fetched all amenities, total count: ", total)
	return amenities, total, err
}

// GetAllWithFilters returns amenities applying filters and pagination.
// Service enforces sane defaults and max page size to centralize pagination logic.
func (s *amenityService) GetAllWithFilters(ctx context.Context, page, size int, amenityType, propertyID string) ([]models.Amenities, int64, error) {
	// Service layer validation - normalize parameters
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}
	amenities, total, err := s.repo.GetAllWithFilters(ctx, page, size, amenityType, propertyID)
	if err != nil {
		return nil, 0, err
	}
	logger.Info("Fetched amenities with filters, total count: ", total)
	return amenities, total, err
}

// GetByID retrieves an amenity by its unique ID.
func (s *amenityService) GetByID(ctx context.Context, id string) (*models.Amenities, error) {
	if id == "" {
		return nil, fmt.Errorf("invalid amenity ID: cannot be empty")
	}
	amenity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched amenity with ID: ", id)
	return amenity, nil
}

// Create adds a new amenity record.
func (s *amenityService) Create(ctx context.Context, amenity *models.Amenities) error {
	if amenity == nil {
		return fmt.Errorf("amenity cannot be nil")
	}
	err := s.repo.Create(ctx, amenity)
	if err != nil {
		return err
	}
	logger.Info("Created amenity with ID: ", amenity.ID)
	return nil
}

// Update modifies an existing amenity record by its ID.
func (s *amenityService) Update(ctx context.Context, id string, amenity *models.Amenities) error {
	if id == "" {
		return fmt.Errorf("invalid amenity ID: cannot be empty")
	}
	if amenity == nil {
		return fmt.Errorf("amenity cannot be nil")
	}
	err := s.repo.Update(ctx, id, amenity)
	if err != nil {
		return err
	}
	logger.Info("Updated amenity with ID: ", id)
	return nil
}

// Delete removes an amenity record by its ID.
func (s *amenityService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("invalid amenity ID: cannot be empty")
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	logger.Info("Deleted amenity with ID: ", id)
	return nil
}

// GetByPropertyID retrieves amenities for a given property ID.
func (s *amenityService) GetByPropertyID(ctx context.Context, propertyID string) (*models.Amenities, error) {
	if propertyID == "" {
		return nil, fmt.Errorf("invalid property ID: cannot be empty")
	}
	amenities, err := s.repo.GetByPropertyID(ctx, propertyID)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched amenities for property ID: ", propertyID)
	return amenities, nil
}
