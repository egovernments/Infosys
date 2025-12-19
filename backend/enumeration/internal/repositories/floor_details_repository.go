package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// floorDetailsRepository implements the FloorDetailsRepository interface using GORM.
var _ FloorDetailsRepository = (*floorDetailsRepository)(nil)

type floorDetailsRepository struct {
	db *gorm.DB // GORM database connection
}

// NewFloorDetailsRepository creates a new repository for floor details.
func NewFloorDetailsRepository(db *gorm.DB) FloorDetailsRepository {
	return &floorDetailsRepository{db: db}
}

// GetAll retrieves all FloorDetails records with optional filtering by construction details ID and supports pagination.
// Returns the filtered floor details and the total count.
func (r *floorDetailsRepository) GetAll(ctx context.Context, page, size int, constructionDetailsID *uuid.UUID) ([]*models.FloorDetails, int64, error) {
	var floorDetails []*models.FloorDetails
	var total int64

	query := r.db.Model(&models.FloorDetails{})

	if constructionDetailsID != nil {
		query = query.Where("construction_details_id = ?", *constructionDetailsID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count floor details: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Find(&floorDetails).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no floor details found")
		}
		return nil, 0, fmt.Errorf("failed to get floor details: %w", err)
	}

	return floorDetails, total, nil
}

// Create inserts a new FloorDetails record into the database.
func (r *floorDetailsRepository) Create(ctx context.Context, floorDetails *models.FloorDetails) error {
	if err := r.db.Create(floorDetails).Error; err != nil {
		return fmt.Errorf("failed to create floor details: %w", err)
	}
	return nil
}

// GetByID retrieves a FloorDetails record by its ID.
func (r *floorDetailsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.FloorDetails, error) {
	var floorDetails models.FloorDetails
	err := r.db.Where("id = ?", id).First(&floorDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("floor details with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get floor details by id %s: %w", id, err)
	}
	return &floorDetails, nil
}

// Update modifies an existing FloorDetails record in the database.
// Returns an error if the record does not exist or update fails.
func (r *floorDetailsRepository) Update(ctx context.Context, floorDetails *models.FloorDetails) error {
	result := r.db.Save(floorDetails)
	if result.Error != nil {
		return fmt.Errorf("failed to update floor details with id %s: %w", floorDetails.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("floor details with id %s not found for update", floorDetails.ID)
	}
	return nil
}

// Delete removes a FloorDetails record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *floorDetailsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.FloorDetails{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete floor details with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("floor details with id %s not found for deletion", id)
	}
	return nil
}

// GetByConstructionDetailsID retrieves all FloorDetails records for a given construction details ID.
func (r *floorDetailsRepository) GetByConstructionDetailsID(ctx context.Context, constructionDetailsID uuid.UUID) ([]*models.FloorDetails, error) {
	var floorDetails []*models.FloorDetails
	err := r.db.Where("construction_details_id = ?", constructionDetailsID).Find(&floorDetails).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no floor details found for construction details id %s", constructionDetailsID)
		}
		return nil, fmt.Errorf("failed to get floor details by construction details id %s: %w", constructionDetailsID, err)
	}
	return floorDetails, nil
}
