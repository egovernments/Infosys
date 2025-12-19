package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

// Compile-time check for interface implementation
var _ FloorDetailsService = (*floorDetailsService)(nil)

// floorDetailsService handles business logic for floor details
type floorDetailsService struct {
	floorDetailsRepo repositories.FloorDetailsRepository
}

// NewFloorDetailsService returns a new floorDetailsService
func NewFloorDetailsService(floorDetailsRepo repositories.FloorDetailsRepository) FloorDetailsService {
	return &floorDetailsService{
		floorDetailsRepo: floorDetailsRepo,
	}
}

// CreateFloorDetails validates and adds a new floor details record
func (s *floorDetailsService) CreateFloorDetails(ctx context.Context, floorDetails *models.FloorDetails) error {
	if floorDetails == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if err := s.validateFloorDetails(floorDetails); err != nil {
		return err
	}

	return s.floorDetailsRepo.Create(ctx, floorDetails)
}

// GetFloorDetailsByID fetches a floor details record by its ID
func (s *floorDetailsService) GetFloorDetailsByID(ctx context.Context, id uuid.UUID) (*models.FloorDetails, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid floor details ID: cannot be nil")
	}

	// Repository layer error - propagate as is
	return s.floorDetailsRepo.GetByID(ctx, id)
}

// UpdateFloorDetails validates and updates an existing floor details record
func (s *floorDetailsService) UpdateFloorDetails(ctx context.Context, floorDetails *models.FloorDetails) error {
	if floorDetails == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if floorDetails.ID == uuid.Nil {
		return fmt.Errorf("invalid floor details ID: cannot be nil")
	}
	if err := s.validateFloorDetails(floorDetails); err != nil {
		return err
	}
	return s.floorDetailsRepo.Update(ctx, floorDetails)
}

// DeleteFloorDetails removes a floor details record by its ID
func (s *floorDetailsService) DeleteFloorDetails(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid floor details ID: cannot be nil")
	}
	return s.floorDetailsRepo.Delete(ctx, id)
}

// GetAllFloorDetails returns all floor details with pagination, optionally filtered by constructionDetailsID
func (s *floorDetailsService) GetAllFloorDetails(ctx context.Context, page, size int, constructionDetailsID *uuid.UUID) ([]*models.FloorDetails, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}

	return s.floorDetailsRepo.GetAll(ctx, page, size, constructionDetailsID)
}

// GetFloorDetailsByConstructionDetailsID fetches all floor details for a construction details ID
func (s *floorDetailsService) GetFloorDetailsByConstructionDetailsID(ctx context.Context, constructionDetailsID uuid.UUID) ([]*models.FloorDetails, error) {
	if constructionDetailsID == uuid.Nil {
		return nil, fmt.Errorf("invalid construction details ID: cannot be nil")
	}

	return s.floorDetailsRepo.GetByConstructionDetailsID(ctx, constructionDetailsID)
}

// validateFloorDetails checks required fields for floor details and returns an error if validation fails
func (s *floorDetailsService) validateFloorDetails(floorDetails *models.FloorDetails) error {
	if floorDetails.ConstructionDetailsID == uuid.Nil {
		return fmt.Errorf("construction details ID is required")
	}

	if floorDetails.LengthFt <= 0 {
		return fmt.Errorf("length in feet must be greater than 0")
	}

	if floorDetails.BreadthFt <= 0 {
		return fmt.Errorf("breadth in feet must be greater than 0")
	}

	if floorDetails.PlinthAreaSqFt <= 0 {
		return fmt.Errorf("plinth area must be greater than 0")
	}

	return nil
}
