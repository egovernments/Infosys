package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Compile-time check for interface implementation
var _ PropertyService = (*propertyService)(nil)

// propertyService handles business logic for properties
type propertyService struct {
	propertyRepo repositories.PropertyRepository
}

// NewPropertyService returns a new propertyService
func NewPropertyService(propertyRepo repositories.PropertyRepository) PropertyService {
	return &propertyService{
		propertyRepo: propertyRepo,
	}
}

// CreateProperty validates and adds a new property, generating a property number if needed
func (s *propertyService) CreateProperty(ctx context.Context, property *models.Property) error {
	if property == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	// Generate property number if not provided
	if property.PropertyNo == "" {
		propertyNo, err := s.GeneratePropertyNo(ctx)
		if err != nil {
			return fmt.Errorf("failed to generate property number: %v", err)
		}
		property.PropertyNo = propertyNo
	}

	// Set PropertyID for nested address if provided
	if property.Address != nil {
		property.Address.PropertyID = property.ID
		if property.Address.ID == uuid.Nil {
			property.Address.ID = uuid.New()
		}
	}

	return s.propertyRepo.Create(ctx, property)
}

// GetPropertyByID fetches a property by its ID
func (s *propertyService) GetPropertyByID(ctx context.Context, id uuid.UUID) (*models.Property, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid property ID: cannot be nil")
	}

	// Repository layer error - propagate as is
	return s.propertyRepo.GetByID(ctx, id)
}

// UpdateProperty validates and updates an existing property
func (s *propertyService) UpdateProperty(ctx context.Context, property *models.Property) error {
	if property == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	if property.ID == uuid.Nil {
		return fmt.Errorf("invalid property ID: cannot be nil")
	}

	if err := s.validateProperty(property); err != nil {
		return err
	}

	// Repository layer error - propagate as is
	return s.propertyRepo.Update(ctx, property)
}

// DeleteProperty removes a property by its ID
func (s *propertyService) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid property ID: cannot be nil")
	}
	// Repository layer error - propagate as is
	return s.propertyRepo.Delete(ctx, id)
}

// GetAllProperties returns all properties with pagination, optionally filtered by propertyType
func (s *propertyService) GetAllProperties(ctx context.Context, page, size int, propertyType *string) ([]*models.Property, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}

	return s.propertyRepo.GetAll(ctx, page, size, propertyType)
}

// GetPropertyByPropertyNo fetches a property by its property number
func (s *propertyService) GetPropertyByPropertyNo(ctx context.Context, propertyNo string) (*models.Property, error) {
	if propertyNo == "" {
		return nil, fmt.Errorf("property number is required")
	}

	// Repository layer error - propagate as is
	return s.propertyRepo.GetByPropertyNo(ctx, propertyNo)
}

// SearchProperties finds properties matching the given search parameters
func (s *propertyService) SearchProperties(ctx context.Context, params SearchPropertyParams) ([]*models.Property, int64, error) {
	if params.Page < 0 {
		params.Page = 0
	}
	if params.Size <= 0 || params.Size > 100 {
		params.Size = 20
	}

	repoParams := repositories.SearchPropertyParams{
		Page:          params.Page,
		Size:          params.Size,
		PropertyType:  params.PropertyType,
		OwnershipType: params.OwnershipType,
		ComplexName:   params.ComplexName,
		Locality:      params.Locality,
		WardNo:        params.WardNo,
		ZoneNo:        params.ZoneNo,
		Street:        params.Street,
		SortBy:        params.SortBy,
		SortOrder:     params.SortOrder,
	}

	return s.propertyRepo.Search(ctx, repoParams)
}

// GeneratePropertyNo creates a new unique property number
func (s *propertyService) GeneratePropertyNo(ctx context.Context) (string, error) {
	// Generate property number in format: PROP-YYYY-XXXXXX
	year := time.Now().Year()
	timestamp := time.Now().Unix()

	propertyNo := fmt.Sprintf("PROP-%d-%06d", year, timestamp%1000000)

	// Check if property number already exists (very unlikely but good to check)
	_, err := s.propertyRepo.GetByPropertyNo(ctx, propertyNo)
	if err == nil {
		// Property number exists, add random suffix
		propertyNo = fmt.Sprintf("%s-%d", propertyNo, time.Now().Nanosecond()%1000)
	}

	return propertyNo, nil
}

// validateProperty checks required fields and allowed values for a property
func (s *propertyService) validateProperty(property *models.Property) error {
	// Validate required fields
	if property.OwnershipType == "" {
		return fmt.Errorf("ownership type is required")
	}

	if property.PropertyType == "" {
		return fmt.Errorf("property type is required")
	}

	
	return nil
}

