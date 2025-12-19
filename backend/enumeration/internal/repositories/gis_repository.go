package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// gisRepository implements the GISRepository interface using GORM.
var _ GISRepository = (*gisRepository)(nil)

type gisRepository struct {
	db *gorm.DB // GORM database connection
}

// NewGISRepository creates a new repository for GIS data.
func NewGISRepository(db *gorm.DB) GISRepository {
	return &gisRepository{db: db}
}

// Create inserts a new GISData record into the database after validating required fields.
// Returns an error if validation fails or GIS data already exists for the property.
func (r *gisRepository) Create(ctx context.Context, gisData *models.GISData) error {
	// Validate required fields
	if gisData.PropertyID == uuid.Nil {
		return fmt.Errorf("property ID is required")
	}
	if gisData.Source == "" {
		return fmt.Errorf("source is required")
	}
	if gisData.Type == "" {
		return fmt.Errorf("type is required")
	}

	// Check if GIS data already exists for this property
	var existingGISData models.GISData
	if err := r.db.Where("property_id = ?", gisData.PropertyID).First(&existingGISData).Error; err == nil {
		return fmt.Errorf("GIS data already exists for property ID %s", gisData.PropertyID)
	}

	// Create the GIS data
	if err := r.db.Create(gisData).Error; err != nil {
		return fmt.Errorf("failed to create GIS data: %w", err)
	}

	return nil
}

// GetByID retrieves a GISData record by its ID, including related Coordinates.
func (r *gisRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.GISData, error) {
	var gisData models.GISData
	if err := r.db.Preload("Coordinates").First(&gisData, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("GIS data not found with ID %s", id)
		}
		return nil, fmt.Errorf("failed to get GIS data: %w", err)
	}
	return &gisData, nil
}

// GetByPropertyID retrieves a GISData record by property ID, including related Coordinates.
func (r *gisRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.GISData, error) {
	var gisData models.GISData
	if err := r.db.Preload("Coordinates").First(&gisData, "property_id = ?", propertyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("GIS data not found for property ID %s", propertyID)
		}
		return nil, fmt.Errorf("failed to get GIS data: %w", err)
	}
	return &gisData, nil
}

// Update modifies an existing GISData record in the database.
// Checks for property ID conflicts and reloads the updated data with related Coordinates.
func (r *gisRepository) Update(ctx context.Context, gisData *models.GISData) error {
	if gisData.ID == uuid.Nil {
		return fmt.Errorf("GIS data ID is required for update")
	}

	// Check if the record exists
	var existingGISData models.GISData
	if err := r.db.First(&existingGISData, "id = ?", gisData.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("GIS data not found with ID %s", gisData.ID)
		}
		return fmt.Errorf("failed to find GIS data: %w", err)
	}

	// If property ID is being changed, check if another GIS data exists for the new property
	if gisData.PropertyID != existingGISData.PropertyID {
		var conflictingGISData models.GISData
		if err := r.db.Where("property_id = ? AND id != ?", gisData.PropertyID, gisData.ID).First(&conflictingGISData).Error; err == nil {
			return fmt.Errorf("GIS data already exists for property ID %s", gisData.PropertyID)
		}
	}

	// Update the record
	if err := r.db.Model(&existingGISData).Updates(gisData).Error; err != nil {
		return fmt.Errorf("failed to update GIS data: %w", err)
	}

	// Reload the updated data
	if err := r.db.Preload("Coordinates").First(gisData, "id = ?", gisData.ID).Error; err != nil {
		return fmt.Errorf("failed to reload updated GIS data: %w", err)
	}

	return nil
}

// Delete removes a GISData record by its ID.
// Returns an error if the record does not exist or deletion fails.
// Related coordinates will be deleted automatically due to CASCADE.
func (r *gisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if the record exists
	var gisData models.GISData
	if err := r.db.First(&gisData, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("GIS data not found with ID %s", id)
		}
		return fmt.Errorf("failed to find GIS data: %w", err)
	}

	// Delete the record (coordinates will be deleted automatically due to CASCADE)
	if err := r.db.Delete(&gisData).Error; err != nil {
		return fmt.Errorf("failed to delete GIS data: %w", err)
	}

	return nil
}

// GetAll retrieves all GISData records with pagination, including related Coordinates.
// Returns the paginated GIS data and the total count.
func (r *gisRepository) GetAll(ctx context.Context, page, size int) ([]*models.GISData, int64, error) {
	var gisDataList []*models.GISData
	var total int64

	// Count total records
	if err := r.db.Model(&models.GISData{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count GIS data: %w", err)
	}

	// Calculate offset
	offset := page * size

	// Retrieve paginated records
	if err := r.db.Preload("Coordinates").
		Offset(offset).
		Limit(size).
		Find(&gisDataList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get GIS data: %w", err)
	}

	return gisDataList, total, nil
}
