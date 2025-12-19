package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// assessmentDetailsRepository implements the AssessmentDetailsRepository interface using GORM.
var _ AssessmentDetailsRepository = (*assessmentDetailsRepository)(nil)

type assessmentDetailsRepository struct {
	db *gorm.DB // GORM database connection
}

// NewAssessmentDetailsRepository creates a new repository for assessment details.
func NewAssessmentDetailsRepository(db *gorm.DB) AssessmentDetailsRepository {
	return &assessmentDetailsRepository{db: db}
}

// Create inserts a new AssessmentDetails record into the database.
func (r *assessmentDetailsRepository) Create(ctx context.Context, assessmentDetails *models.AssessmentDetails) error {
	res := r.db.Create(assessmentDetails)
	if res.Error != nil {
		return fmt.Errorf("failed to create assessment details: %w", res.Error)
	}
	return nil
}

// GetByID retrieves an AssessmentDetails record by its ID.
func (r *assessmentDetailsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AssessmentDetails, error) {
	var assessmentDetails models.AssessmentDetails
	err := r.db.Where("id = ?", id).First(&assessmentDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assessment details with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get assessment details by id %s: %w", id, err)
	}
	return &assessmentDetails, nil
}

// Update modifies an existing AssessmentDetails record in the database.
// Returns an error if the record does not exist or update fails.
func (r *assessmentDetailsRepository) Update(ctx context.Context, assessmentDetails *models.AssessmentDetails) error {
	result := r.db.Save(assessmentDetails)
	if result.Error != nil {
		return fmt.Errorf("failed to update assessment details with id %s: %w", assessmentDetails.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("assessment details with id %s not found for update", assessmentDetails.ID)
	}
	return nil
}

// Delete removes an AssessmentDetails record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *assessmentDetailsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.AssessmentDetails{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete assessment details with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("assessment details with id %s not found for deletion", id)
	}
	return nil
}

// GetAll retrieves all AssessmentDetails records with optional filtering by property ID and supports pagination.
// Returns the filtered assessment details and the total count.
func (r *assessmentDetailsRepository) GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.AssessmentDetails, int64, error) {
	var assessmentDetails []*models.AssessmentDetails
	var total int64

	query := r.db.Model(&models.AssessmentDetails{})

	if propertyID != nil {
		query = query.Where("property_id = ?", *propertyID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count assessment details: %w", err)
	}

	// Apply pagination
	offset := page * size
	err := query.Offset(offset).Limit(size).Find(&assessmentDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no assessment details found")
		}
		return nil, 0, fmt.Errorf("failed to get assessment details: %w", err)
	}

	return assessmentDetails, total, nil
}

// GetByPropertyID retrieves an AssessmentDetails record by the associated property ID.
func (r *assessmentDetailsRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.AssessmentDetails, error) {
	var assessmentDetails models.AssessmentDetails
	err := r.db.Where("property_id = ?", propertyID).First(&assessmentDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assessment details with property id %s not found", propertyID)
		}
		return nil, fmt.Errorf("failed to get assessment details by property id %s: %w", propertyID, err)
	}
	return &assessmentDetails, nil
}
