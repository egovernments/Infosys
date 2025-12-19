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
var _ GISService = (*gisService)(nil)

// gisService handles business logic for GIS data
type gisService struct {
	gisRepo repositories.GISRepository
}

// NewGISService returns a new gisService
func NewGISService(gisRepo repositories.GISRepository) GISService {
	return &gisService{
		gisRepo: gisRepo,
	}
}

// CreateGISData validates and adds a new GIS data record
func (s *gisService) CreateGISData(ctx context.Context, gisData *models.GISData) error {
	// Validate input
	if gisData == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if err := s.validateGISData(gisData); err != nil {
		return err
	}
	return s.gisRepo.Create(ctx, gisData)
}

// GetGISDataByID fetches a GIS data record by its ID
func (s *gisService) GetGISDataByID(ctx context.Context, id uuid.UUID) (*models.GISData, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("GIS data ID cannot be empty")
	}
	return s.gisRepo.GetByID(ctx, id)
}

// GetGISDataByPropertyID fetches GIS data for a property
func (s *gisService) GetGISDataByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.GISData, error) {
	if propertyID == uuid.Nil {
		return nil, fmt.Errorf("property ID cannot be empty")
	}
	return s.gisRepo.GetByPropertyID(ctx, propertyID)
}

// UpdateGISData validates and updates an existing GIS data record
func (s *gisService) UpdateGISData(ctx context.Context, gisData *models.GISData) error {
	if gisData == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	if gisData.ID == uuid.Nil {
		return fmt.Errorf("GIS data ID cannot be empty for update")
	}

	// Validate input
	if err := s.validateGISData(gisData); err != nil {
		return err
	}
	return s.gisRepo.Update(ctx, gisData)
}

// DeleteGISData removes a GIS data record by its ID
func (s *gisService) DeleteGISData(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("GIS data ID cannot be empty")
	}
	return s.gisRepo.Delete(ctx, id)
}

// GetAllGISData returns all GIS data with pagination
func (s *gisService) GetAllGISData(ctx context.Context, page, size int) ([]*models.GISData, int64, error) {
	// Validate pagination parameters
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}
	//propagate the repository error as is
	return s.gisRepo.GetAll(ctx, page, size)
}

// validateGISData checks required fields and allowed values for GIS data
func (s *gisService) validateGISData(gisData *models.GISData) error {
	if gisData == nil {
		return fmt.Errorf("GIS data cannot be nil")
	}

	if gisData.PropertyID == uuid.Nil {
		return fmt.Errorf("property ID is required")
	}

	if gisData.Source == "" {
		return fmt.Errorf("source is required")
	}

	// Validate source enum values
	validSources := map[string]bool{
		"GPS":          true,
		"MANUAL_ENTRY": true,
		"IMPORT":       true,
	}
	if !validSources[gisData.Source] {
		return fmt.Errorf("invalid source: %s. Valid values are: GPS, MANUAL_ENTRY, IMPORT", gisData.Source)
	}

	if gisData.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Validate type enum values
	validTypes := map[string]bool{
		"POINT":   true,
		"LINE":    true,
		"POLYGON": true,
	}
	if !validTypes[gisData.Type] {
		return fmt.Errorf("invalid type: %s. Valid values are: POINT, LINE, POLYGON", gisData.Type)
	}

	// Entity type is optional, but if provided, should not be empty
	if gisData.EntityType != "" && len(gisData.EntityType) > 100 {
		return fmt.Errorf("entity type cannot exceed 100 characters")
	}

	return nil
}
