package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// igrsRepository implements the IGRSRepository interface using GORM for database operations.
var _ IGRSRepository = (*igrsRepository)(nil)

type igrsRepository struct {
	db *gorm.DB // GORM database connection
}

// NewIGRSRepository creates a new repository for IGRS data.
func NewIGRSRepository(db *gorm.DB) IGRSRepository {
	return &igrsRepository{db: db}
}

// Create inserts a new IGRS record into the database.
// Returns an error if the operation fails.
func (r *igrsRepository) Create(ctx context.Context, igrs *models.IGRS) error {
	if err := r.db.WithContext(ctx).Create(igrs).Error; err != nil {
		return fmt.Errorf("failed to create igrs: %w", err)
	}
	return nil
}

// GetByPropertyID retrieves an IGRS record by property ID.
// Returns the IGRS record or an error if not found.
func (r *igrsRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID) (*models.IGRS, error) {
	var m models.IGRS
	if err := r.db.WithContext(ctx).First(&m, "property_id = ?", propertyID).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByID retrieves an IGRS record by its ID.
// Returns the IGRS record or an error if not found.
func (r *igrsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.IGRS, error) {
	var m models.IGRS
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Update modifies an existing IGRS record in the database.
// Returns an error if the update fails.
func (r *igrsRepository) Update(ctx context.Context, igrs *models.IGRS) error {
	if err := r.db.WithContext(ctx).Save(igrs).Error; err != nil {
		return fmt.Errorf("failed to update igrs: %w", err)
	}
	return nil
}

// Delete removes an IGRS record by its ID.
// Returns an error if the deletion fails.
func (r *igrsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.IGRS{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete igrs: %w", err)
	}
	return nil
}

// FindAll retrieves all IGRS records with pagination.
// Returns the list of IGRS records, total count, or an error.
func (r *igrsRepository) FindAll(ctx context.Context, page, size int) ([]models.IGRS, int64, error) {
	var list []models.IGRS
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.IGRS{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := page * size
	if err := r.db.WithContext(ctx).Limit(size).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
