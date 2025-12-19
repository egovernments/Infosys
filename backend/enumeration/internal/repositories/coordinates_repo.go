package repositories

import (
	"context"
	"enumeration/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// coordinatesRepository implements the CoordinatesRepository interface using GORM.
var _ CoordinatesRepository = (*coordinatesRepository)(nil)

// coordinatesRepository handles database operations for coordinates.
type coordinatesRepository struct {
	db *gorm.DB // GORM database connection
}

// NewCoordinatesRepository creates a new repository for coordinates.
func NewCoordinatesRepository(db *gorm.DB) CoordinatesRepository {
	return &coordinatesRepository{db: db}
}

// Create inserts a new Coordinates record into the database.
func (r *coordinatesRepository) Create(ctx context.Context, coordinates *models.Coordinates) error {
	if err := r.db.Create(coordinates).Error; err != nil {
		return fmt.Errorf("failed to create coordinates: %w", err)
	}
	return nil
}

// GetByID retrieves a Coordinates record by its ID.
func (r *coordinatesRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Coordinates, error) {
	var coordinates models.Coordinates
	err := r.db.First(&coordinates, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("coordinates with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get coordinates by id %s: %w", id, err)
	}
	return &coordinates, nil
}

// Update modifies an existing Coordinates record in the database.
func (r *coordinatesRepository) Update(ctx context.Context, coordinates *models.Coordinates) error {
	if err := r.db.Save(coordinates).Error; err != nil {
		return fmt.Errorf("failed to update coordinates with id %s: %w", coordinates.ID, err)
	}
	return nil
}

// Delete removes a Coordinates record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *coordinatesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.Coordinates{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete coordinates with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("coordinates with id %s not found for deletion", id)
	}
	return nil
}

// FindAll retrieves all Coordinates records with optional filtering by GIS Data ID and supports pagination.
// Returns the filtered coordinates and the total count.
func (r *coordinatesRepository) FindAll(ctx context.Context, page, size int, gisDataID *uuid.UUID) ([]models.Coordinates, int64, error) {
	var coordinates []models.Coordinates
	var total int64

	query := r.db.Model(&models.Coordinates{})

	// Filter by GIS Data ID if provided
	if gisDataID != nil {
		query = query.Where("gis_data_id = ?", *gisDataID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count coordinates: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Find(&coordinates).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch coordinates with pagination: %w", err)
	}

	return coordinates, total, nil
}

// FindByGISDataID retrieves all Coordinates records for a specific GIS data record.
func (r *coordinatesRepository) FindByGISDataID(ctx context.Context, gisDataID uuid.UUID) ([]models.Coordinates, error) {
	var coordinates []models.Coordinates
	if err := r.db.Where("gis_data_id = ?", gisDataID).Find(&coordinates).Error; err != nil {
		return nil, fmt.Errorf("failed to get coordinates by gis data id %s: %w", gisDataID, err)
	}
	return coordinates, nil
}

// DeleteByGISDataID deletes all Coordinates records for a specific GIS data record.
func (r *coordinatesRepository) DeleteByGISDataID(ctx context.Context, gisDataID uuid.UUID) error {
	if err := r.db.Where("gis_data_id = ?", gisDataID).Delete(&models.Coordinates{}).Error; err != nil {
		return fmt.Errorf("failed to delete coordinates by gis data id %s: %w", gisDataID, err)
	}
	return nil
}

// CreateBatch inserts multiple Coordinates records in a single transaction for efficiency and rollback on failure.
func (r *coordinatesRepository) CreateBatch(ctx context.Context, coords []*models.Coordinates) error {
	// Use the DB context, transaction and CreateInBatches for efficiency and rollback on failure
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(coords, 100).Error; err != nil {
			return fmt.Errorf("failed to create coordinates in batch: %w", err)
		}
		return nil
	})
}

// ReplaceByGISDataID replaces all Coordinates records for a specific GIS data record with a new batch.
// Deletes existing records and inserts the new batch in a single transaction.
func (r *coordinatesRepository) ReplaceByGISDataID(ctx context.Context, gisDataID uuid.UUID, coords []*models.Coordinates) error {
	// Transaction: delete existing -> create new batch
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing coordinates for the GIS data
		if err := tx.Where("gis_data_id = ?", gisDataID).Delete(&models.Coordinates{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing coordinates for gis data id %s: %w", gisDataID, err)
		}

		// Ensure each coordinate has the correct GISDataID and an ID
		for _, c := range coords {
			if c == nil {
				continue
			}
			c.GISDataID = gisDataID
			if c.ID == uuid.Nil {
				c.ID = uuid.New()
			}
		}

		// If nothing to insert, that's fine (we've deleted existing)
		if len(coords) == 0 {
			return nil
		}

		// Insert new coordinates in batches
		if err := tx.CreateInBatches(coords, 100).Error; err != nil {
			return fmt.Errorf("failed to create coordinates in batch: %w", err)
		}
		return nil
	})
}
