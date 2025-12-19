package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ensure propertyOwnerRepository implements PropertyOwnerRepository interface at compile time.
var _ PropertyOwnerRepository = (*propertyOwnerRepository)(nil)

// propertyOwnerRepository provides implementation for PropertyOwnerRepository using GORM for database operations.
type propertyOwnerRepository struct {
	db *gorm.DB
}

// NewPropertyOwnerRepository creates a new instance of propertyOwnerRepository.
// db: GORM database connection.
// Returns: PropertyOwnerRepository implementation.
func NewPropertyOwnerRepository(db *gorm.DB) PropertyOwnerRepository {
	return &propertyOwnerRepository{db: db}
}

// CreateBatch inserts multiple PropertyOwner records in batches using a transaction for efficiency.
// ctx: context for the operation.
// owners: slice of PropertyOwner pointers to be created.
// Returns: error if batch creation fails.
func (r *propertyOwnerRepository) CreateBatch(ctx context.Context, owners []*models.PropertyOwner) error {
	// Use a transaction and CreateInBatches for efficiency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(owners) == 0 {
			return nil
		}
		// Adjust batch size if needed
		if err := tx.CreateInBatches(owners, 100).Error; err != nil {
			return fmt.Errorf("failed to create property owners batch: %w", err)
		}
		return nil
	})
}

// Create inserts a new PropertyOwner record into the database.
// ctx: context for the operation.
// owner: pointer to PropertyOwner model to be created.
// Returns: error if creation fails.
func (r *propertyOwnerRepository) Create(ctx context.Context, owner *models.PropertyOwner) error {

	if err := r.db.WithContext(ctx).Create(owner).Error; err != nil {
		return fmt.Errorf("failed to create property owner: %w", err)
	}
	return nil
}

// GetAll retrieves all PropertyOwner records from the database.
// ctx: context for the operation.
// Returns: slice of PropertyOwner pointers and error if any.
func (r *propertyOwnerRepository) GetAll(ctx context.Context) ([]*models.PropertyOwner, error) {
	var owners []*models.PropertyOwner
	err := r.db.WithContext(ctx).Find(&owners).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no property owners found")
		}
		return nil, fmt.Errorf("failed to get property owners: %w", err)
	}
	return owners, nil
}

// GetByID retrieves a PropertyOwner by its unique ID.
// ctx: context for the operation.
// id: UUID of the property owner.
// Returns: pointer to PropertyOwner and error if not found or on failure.
func (r *propertyOwnerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PropertyOwner, error) {
	var owner models.PropertyOwner
	err := r.db.WithContext(ctx).First(&owner, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("property owner with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get property owner by id %s: %w", id, err)
	}
	return &owner, nil
}

// GetByPropertyID retrieves property owners for a property, with manual pagination handled in the repository.
// ctx: context for the operation.
// propertyID: UUID of the property.
// page: page number (zero-based).
// size: number of records per page.
// Returns: slice of PropertyOwner pointers, total count, and error if any.
func (r *propertyOwnerRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID, page, size int) ([]*models.PropertyOwner, int64, error) {
	// clamp pagination
	if page < 0 {
		page = 0
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	var total int64
	query := r.db.WithContext(ctx).Model(&models.PropertyOwner{}).Where("property_id = ?", propertyID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count property owners for property id %s: %w", propertyID, err)
	}

	// if no rows, return empty slice + total 0
	if total == 0 {
		return []*models.PropertyOwner{}, 0, nil
	}

	var owners []*models.PropertyOwner
	offset := page * size
	if err := query.Order("created_at asc").Limit(size).Offset(offset).Find(&owners).Error; err != nil {
		// treat no rows as empty result
		if err == gorm.ErrRecordNotFound {
			return []*models.PropertyOwner{}, total, nil
		}
		return nil, 0, fmt.Errorf("failed to get property owners by property id %s: %w", propertyID, err)
	}

	return owners, total, nil
}

// Update modifies an existing PropertyOwner record in the database.
// ctx: context for the operation.
// owner: pointer to PropertyOwner model with updated data.
// Returns: error if update fails or record not found.
func (r *propertyOwnerRepository) Update(ctx context.Context, owner *models.PropertyOwner) error {
	result := r.db.WithContext(ctx).Save(owner)
	if result.Error != nil {
		return fmt.Errorf("failed to update property owner with id %s: %w", owner.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("property owner with id %s not found for update", owner.ID)
	}
	return nil
}

// Delete removes a PropertyOwner record by its unique ID.
// ctx: context for the operation.
// id: UUID of the property owner to delete.
// Returns: error if deletion fails or record not found.
func (r *propertyOwnerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.PropertyOwner{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete property owner with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("property owner with id %s not found for deletion", id)
	}
	return nil
}
