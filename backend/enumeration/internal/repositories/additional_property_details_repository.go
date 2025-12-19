package repositories

import (
	"context"
	"enumeration/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// additionalPropertyDetailsRepository implements the AdditionalPropertyDetailsRepository interface using GORM.
var _ AdditionalPropertyDetailsRepository = (*additionalPropertyDetailsRepository)(nil)

type additionalPropertyDetailsRepository struct {
	db *gorm.DB // GORM database connection
}

// NewAdditionalPropertyDetailsRepository creates a new repository for additional property details.
func NewAdditionalPropertyDetailsRepository(db *gorm.DB) AdditionalPropertyDetailsRepository {
	return &additionalPropertyDetailsRepository{db: db}
}

// Create inserts a new AdditionalPropertyDetails record into the database.
func (r *additionalPropertyDetailsRepository) Create(ctx context.Context, details *models.AdditionalPropertyDetails) error {
	if err := r.db.Create(details).Error; err != nil {
		return fmt.Errorf("failed to create additional property details: %w", err)
	}
	return nil
}

// GetByID retrieves an AdditionalPropertyDetails record by its ID.
func (r *additionalPropertyDetailsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AdditionalPropertyDetails, error) {
	var details models.AdditionalPropertyDetails
	err := r.db.Where("id = ?", id).First(&details).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("additional property details with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get additional property details by id %s: %w", id, err)
	}
	return &details, nil
}

// Update modifies an existing AdditionalPropertyDetails record in the database.
func (r *additionalPropertyDetailsRepository) Update(ctx context.Context, details *models.AdditionalPropertyDetails) error {
	if err := r.db.Save(details).Error; err != nil {
		return fmt.Errorf("failed to update additional property details with id %s: %w", details.ID, err)
	}
	return nil
}

// Delete removes an AdditionalPropertyDetails record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *additionalPropertyDetailsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.AdditionalPropertyDetails{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete additional property details with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("additional property details with id %s not found for deletion", id)
	}
	return nil
}

// GetAll retrieves all AdditionalPropertyDetails records with optional filtering and pagination.
// Supports filtering by property ID and field name, and returns the total count.
func (r *additionalPropertyDetailsRepository) GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID, fieldName *string) ([]*models.AdditionalPropertyDetails, int64, error) {
	var details []*models.AdditionalPropertyDetails
	var total int64

	query := r.db.Model(&models.AdditionalPropertyDetails{})

	if propertyID != nil {
		query = query.Where("property_id = ?", *propertyID)
	}

	if fieldName != nil && *fieldName != "" {
		query = query.Where("field_name ILIKE ?", "%"+*fieldName+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count additional property details: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Find(&details).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch additional property details with pagination: %w", err)
	}

	return details, total, nil
}

// GetByPropertyID retrieves all AdditionalPropertyDetails records for a given property ID.
func (r *additionalPropertyDetailsRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.AdditionalPropertyDetails, error) {
	var details []*models.AdditionalPropertyDetails
	if err := r.db.Where("property_id = ?", propertyID).Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to get additional property details by property id %s: %w", propertyID, err)
	}
	return details, nil
}

// GetByFieldName retrieves all AdditionalPropertyDetails records for a given field name.
func (r *additionalPropertyDetailsRepository) GetByFieldName(ctx context.Context, fieldName string) ([]*models.AdditionalPropertyDetails, error) {
	var details []*models.AdditionalPropertyDetails
	if err := r.db.Where("field_name = ?", fieldName).Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to get additional property details by field name %s: %w", fieldName, err)
	}
	return details, nil
}
