// Package services provides business logic and service layer implementations
// for the property tax enumeration system.
package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

// Compile-time check to ensure assessmentDetailsService implements AssessmentDetailsService interface.
var _ AssessmentDetailsService = (*assessmentDetailsService)(nil)

// assessmentDetailsService implements the AssessmentDetailsService interface and provides
// methods for managing assessment details records.
type assessmentDetailsService struct {
	assessmentDetailsRepo repositories.AssessmentDetailsRepository // Repository for assessment details data access
}

// NewAssessmentDetailsService creates a new instance of AssessmentDetailsService with the provided repository.
func NewAssessmentDetailsService(assessmentDetailsRepo repositories.AssessmentDetailsRepository) AssessmentDetailsService {
	return &assessmentDetailsService{
		assessmentDetailsRepo: assessmentDetailsRepo,
	}
}

// CreateAssessmentDetails validates and creates a new assessment details record.
func (s *assessmentDetailsService) CreateAssessmentDetails(ctx context.Context, assessmentDetails *models.AssessmentDetails) error {
	if assessmentDetails == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if err := s.validateAssessmentDetails(assessmentDetails); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return s.assessmentDetailsRepo.Create(ctx, assessmentDetails)
}

// GetAssessmentDetailsByID retrieves an assessment details record by its unique ID.
func (s *assessmentDetailsService) GetAssessmentDetailsByID(ctx context.Context, id uuid.UUID) (*models.AssessmentDetails, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid assessment details ID: cannot be nil")
	}
	return s.assessmentDetailsRepo.GetByID(ctx, id)
}

// UpdateAssessmentDetails validates and updates an existing assessment details record.
func (s *assessmentDetailsService) UpdateAssessmentDetails(ctx context.Context, assessmentDetails *models.AssessmentDetails) error {
	if assessmentDetails == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	if assessmentDetails.ID == uuid.Nil {
		return fmt.Errorf("invalid assessment details ID: cannot be nil")
	}

	if err := s.validateAssessmentDetails(assessmentDetails); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return s.assessmentDetailsRepo.Update(ctx, assessmentDetails)
}

// DeleteAssessmentDetails removes an assessment details record by its ID.
func (s *assessmentDetailsService) DeleteAssessmentDetails(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("invalid assessment details ID: cannot be nil")
	}
	return s.assessmentDetailsRepo.Delete(ctx, id)
}

// GetAllAssessmentDetails returns a paginated list of assessment details records and the total count.
// If page or size are invalid, defaults are used. Optionally filters by property ID.
func (s *assessmentDetailsService) GetAllAssessmentDetails(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.AssessmentDetails, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}
	return s.assessmentDetailsRepo.GetAll(ctx, page, size, propertyID)
}

// GetAssessmentDetailsByPropertyID retrieves an assessment details record by property ID.
func (s *assessmentDetailsService) GetAssessmentDetailsByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.AssessmentDetails, error) {
	if propertyID == uuid.Nil {
		return nil, fmt.Errorf("invalid property ID: cannot be nil")
	}
	return s.assessmentDetailsRepo.GetByPropertyID(ctx, propertyID)
}

// validateAssessmentDetails checks the validity of assessment details fields.
// Returns an error if any required field is missing or invalid.
func (s *assessmentDetailsService) validateAssessmentDetails(assessmentDetails *models.AssessmentDetails) error {
	if assessmentDetails.PropertyID == uuid.Nil {
		return fmt.Errorf("property ID is required: cannot be nil")
	}

	// Validate required fields according to OpenAPI spec
	if assessmentDetails.ExtentOfSite == "" {
		return fmt.Errorf("extent of site is required: cannot be empty")
	}

	// Add more validation as needed based on your business rules
	if assessmentDetails.OccupancyCertificateNumber != "" && assessmentDetails.OccupancyCertificateDate == nil {
		return fmt.Errorf("occupancy certificate date is required when certificate number is provided")
	}

	return nil
}
