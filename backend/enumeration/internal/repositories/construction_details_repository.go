package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// constructionDetailsRepository implements the ConstructionDetailsRepository interface using GORM.
var _ ConstructionDetailsRepository = (*constructionDetailsRepository)(nil)

type constructionDetailsRepository struct {
	db *gorm.DB // GORM database connection
}

// NewConstructionDetailsRepository creates a new repository for construction details.
func NewConstructionDetailsRepository(db *gorm.DB) ConstructionDetailsRepository {
	return &constructionDetailsRepository{db: db}
}

// Create inserts a new ConstructionDetails record into the database.
func (r *constructionDetailsRepository) Create(ctx context.Context, constructionDetails *models.ConstructionDetails) error {
	if err := r.db.Create(constructionDetails).Error; err != nil {
		return fmt.Errorf("failed to create construction details: %w", err)
	}
	return nil
}

// GetByID retrieves a ConstructionDetails record by its ID, including related FloorDetails.
func (r *constructionDetailsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ConstructionDetails, error) {
	var constructionDetails models.ConstructionDetails
	err := r.db.Preload("FloorDetails").Where("id = ?", id).First(&constructionDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("construction details with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get construction details by id %s: %w", id, err)
	}
	return &constructionDetails, nil
}

// Update modifies an existing ConstructionDetails record in the database.
// Returns an error if the record does not exist or update fails.
func (r *constructionDetailsRepository) Update(ctx context.Context, constructionDetails *models.ConstructionDetails) error {
	result := r.db.Save(constructionDetails)
	if result.Error != nil {
		return fmt.Errorf("failed to update construction details with id %s: %w", constructionDetails.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("construction details with id %s not found for update", constructionDetails.ID)
	}
	return nil
}

// Delete removes a ConstructionDetails record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *constructionDetailsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.ConstructionDetails{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete construction details with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("construction details with id %s not found for deletion", id)
	}
	return nil
}

// GetAll retrieves all ConstructionDetails records with optional filtering by property ID and supports pagination.
// Preloads related FloorDetails and returns the filtered construction details and the total count.
func (r *constructionDetailsRepository) GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.ConstructionDetails, int64, error) {
	var constructionDetails []*models.ConstructionDetails
	var total int64

	query := r.db.Model(&models.ConstructionDetails{}).Preload("FloorDetails")

	if propertyID != nil {
		query = query.Where("property_id = ?", *propertyID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count construction details: %w", err)
	}

	// Apply pagination
	offset := page * size
	err := query.Offset(offset).Limit(size).Find(&constructionDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no construction details found")
		}
		return nil, 0, fmt.Errorf("failed to get construction details: %w", err)
	}

	return constructionDetails, total, nil
}

// GetByPropertyID retrieves all ConstructionDetails records for a given property ID, including related FloorDetails.
func (r *constructionDetailsRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.ConstructionDetails, error) {
	var constructionDetails []*models.ConstructionDetails
	err := r.db.Preload("FloorDetails").Where("property_id = ?", propertyID).Find(&constructionDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no construction details found for property id %s", propertyID)
		}
		return nil, fmt.Errorf("failed to get construction details by property id %s: %w", propertyID, err)
	}
	return constructionDetails, nil
}
