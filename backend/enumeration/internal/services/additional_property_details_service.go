// Package services provides business logic and service layer implementations
// for the property tax enumeration system.
package services

import (
	"context"
	"encoding/json"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"enumeration/pkg/logger"
	"fmt"

	"github.com/google/uuid"
)

// Compile-time check to ensure additionalPropertyDetailsService implements AdditionalPropertyDetailsService interface.
var _ AdditionalPropertyDetailsService = (*additionalPropertyDetailsService)(nil)

// additionalPropertyDetailsService implements the AdditionalPropertyDetailsService interface and provides
// methods for managing additional property details records.
type additionalPropertyDetailsService struct {
	repo repositories.AdditionalPropertyDetailsRepository // Repository for additional property details data access
}

// NewAdditionalPropertyDetailsService creates a new instance of AdditionalPropertyDetailsService with the provided repository.
func NewAdditionalPropertyDetailsService(repo repositories.AdditionalPropertyDetailsRepository) AdditionalPropertyDetailsService {
	return &additionalPropertyDetailsService{
		repo: repo,
	}
}

// CreateAdditionalPropertyDetails validates and adds a new additional property details record.
func (s *additionalPropertyDetailsService) CreateAdditionalPropertyDetails(ctx context.Context, details *models.AdditionalPropertyDetails) error {
	if details == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	if err := s.validateAdditionalPropertyDetails(details); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	err := s.repo.Create(ctx, details)
	if err != nil {
		return err
	}
	logger.Info("Created additional property details with ID: ", details.ID)
	return nil
}

// GetAdditionalPropertyDetailsByID fetches an additional property details record by its unique ID.
func (s *additionalPropertyDetailsService) GetAdditionalPropertyDetailsByID(ctx context.Context, id uuid.UUID) (*models.AdditionalPropertyDetails, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid additional property details ID: cannot be nil")
	}
	details, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched additional property details with ID: ", id)
	return details, nil
}

// UpdateAdditionalPropertyDetails validates and updates an existing additional property details record.
func (s *additionalPropertyDetailsService) UpdateAdditionalPropertyDetails(ctx context.Context, details *models.AdditionalPropertyDetails) error {
	if details == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	if details.ID == uuid.Nil {
		return fmt.Errorf("invalid additional property details ID: cannot be nil")
	}
	if err := s.validateAdditionalPropertyDetails(details); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	err := s.repo.Update(ctx, details)
	if err != nil {
		return err
	}
	logger.Info("Updated additional property details with ID: ", details.ID)
	return nil
}

// DeleteAdditionalPropertyDetails removes an additional property details record by its ID.
func (s *additionalPropertyDetailsService) DeleteAdditionalPropertyDetails(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid additional property details ID: cannot be nil")
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	logger.Info("Deleted additional property details with ID: ", id)
	return nil
}

// GetAllAdditionalPropertyDetails returns a paginated list of additional property details records and the total count.
// If page or size are invalid, defaults are used. Optionally filters by property ID and field name.
func (s *additionalPropertyDetailsService) GetAllAdditionalPropertyDetails(ctx context.Context, page, size int, propertyID *uuid.UUID, fieldName *string) ([]*models.AdditionalPropertyDetails, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}

	details, total, err := s.repo.GetAll(ctx, page, size, propertyID, fieldName)
	if err != nil {
		return nil, 0, err
	}
	logger.Info("Fetched all additional property details, total count: ", total)
	return details, total, nil
}

// GetAdditionalPropertyDetailsByPropertyID fetches all additional property details records for a given property ID.
func (s *additionalPropertyDetailsService) GetAdditionalPropertyDetailsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.AdditionalPropertyDetails, error) {
	if propertyID == uuid.Nil {
		return nil, fmt.Errorf("invalid property ID: cannot be nil")
	}

	details, err := s.repo.GetByPropertyID(ctx, propertyID)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched additional property details for property ID: ", propertyID)
	return details, nil
}

// GetAdditionalPropertyDetailsByFieldName fetches all additional property details records for a given field name.
func (s *additionalPropertyDetailsService) GetAdditionalPropertyDetailsByFieldName(ctx context.Context, fieldName string) ([]*models.AdditionalPropertyDetails, error) {
	if fieldName == "" {
		return nil, fmt.Errorf("field name is required: cannot be empty")
	}

	details, err := s.repo.GetByFieldName(ctx, fieldName)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched additional property details for field name: ", fieldName)
	return details, nil
}

// validateAdditionalPropertyDetails checks the validity of additional property details fields.
// Returns an error if any required field is missing or invalid, or if FieldValue is not valid JSON.
func (s *additionalPropertyDetailsService) validateAdditionalPropertyDetails(details *models.AdditionalPropertyDetails) error {
	if details.PropertyID == uuid.Nil {
		return fmt.Errorf("property ID is required: cannot be nil")
	}
	if details.FieldName == "" {
		return fmt.Errorf("field name is required: cannot be empty")
	}
	if len(details.FieldValue) == 0 {
		return fmt.Errorf("field value is required: cannot be empty")
	}
	// Validate that FieldValue is valid JSON
	var jsonData interface{}
	if err := json.Unmarshal(details.FieldValue, &jsonData); err != nil {
		return fmt.Errorf("field value must be valid JSON: %w", err)
	}
	return nil
}
