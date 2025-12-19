package repositories

import (
	"context"
	"enumeration/internal/models"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// amenityRepository implements the AmenityRepository interface using GORM.
var _ AmenityRepository = (*amenityRepository)(nil)

type amenityRepository struct {
	db *gorm.DB // GORM database connection
}

// NewAmenityRepository creates a new repository for amenities.
func NewAmenityRepository(db *gorm.DB) AmenityRepository {
	return &amenityRepository{db}
}

// GetAll retrieves all amenities from the database.
func (r *amenityRepository) GetAll(ctx context.Context) ([]models.Amenities, error) {
	var amenities []models.Amenities
	if err := r.db.Find(&amenities).Error; err != nil {
		return nil, fmt.Errorf("failed to get all amenities: %w", err)
	}
	return amenities, nil
}

// GetAllWithFilters retrieves amenities with optional filtering by type and property ID, and supports pagination.
// Returns the filtered amenities and the total count.
func (r *amenityRepository) GetAllWithFilters(ctx context.Context, page, size int, amenityType, propertyID string) ([]models.Amenities, int64, error) {
	var amenities []models.Amenities
	var total int64

	query := r.db.Model(&models.Amenities{})

	// Apply filters
	if amenityType != "" {
		query = query.Where("type = ?", amenityType)
	}
	if propertyID != "" {
		query = query.Where("property_id = ?", propertyID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count amenities: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Limit(size).Offset(offset).Find(&amenities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch amenities with pagination: %w", err)
	}

	return amenities, total, nil
}

// GetByID retrieves an amenity by its ID.
func (r *amenityRepository) GetByID(ctx context.Context, id string) (*models.Amenities, error) {
	var amenity models.Amenities
	err := r.db.First(&amenity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("amenity with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get amenity by id %s: %w", id, err)
	}
	return &amenity, nil
}

// Create inserts a new amenity record into the database.
func (r *amenityRepository) Create(ctx context.Context, amenity *models.Amenities) error {
	if err := r.db.Create(amenity).Error; err != nil {
		return fmt.Errorf("failed to create amenity: %w", err)
	}
	return nil
}

// Update modifies an existing amenity record by its ID.
// Returns an error if the record does not exist or update fails.
func (r *amenityRepository) Update(ctx context.Context, id string, amenity *models.Amenities) error {
	result := r.db.Model(&models.Amenities{}).Where("id = ?", id).Updates(amenity)
	if result.Error != nil {
		return fmt.Errorf("failed to update amenity with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("amenity with id %s not found for update", id)
	}
	return nil
}

// Delete removes an amenity record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *amenityRepository) Delete(ctx context.Context, id string) error {
	result := r.db.Delete(&models.Amenities{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete amenity with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("amenity with id %s not found for deletion", id)
	}
	return nil
}

// GetByPropertyID retrieves an amenity by the associated property ID.
// Returns an error if not found or on failure.
func (r *amenityRepository) GetByPropertyID(ctx context.Context, propertyID string) (*models.Amenities, error) {
	var amenity models.Amenities                                                            // Changed from pointer to value
	err := r.db.WithContext(ctx).Where("property_id = ?", propertyID).First(&amenity).Error // Use First() instead of Find()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("amenity with property ID %s not found", propertyID)
		}
		log.Println("Error fetching amenity by property ID:", err)
		return nil, fmt.Errorf("failed to get amenity by property ID %s: %w", propertyID, err)
	}
	return &amenity, nil
}
