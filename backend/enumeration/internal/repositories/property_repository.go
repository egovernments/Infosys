package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ensure propertyRepository implements PropertyRepository interface at compile time.
var _ PropertyRepository = (*propertyRepository)(nil)

// propertyRepository provides implementation for PropertyRepository using GORM for database operations.
type propertyRepository struct {
	db *gorm.DB
}

// NewPropertyRepository creates a new instance of propertyRepository.
// db: GORM database connection.
// Returns: PropertyRepository implementation.
func NewPropertyRepository(db *gorm.DB) PropertyRepository {
	return &propertyRepository{db: db}
}

// Create inserts a new Property record and its address (if provided) into the database using a transaction.
// ctx: context for the operation.
// property: pointer to Property model to be created.
// Returns: error if creation fails.
func (r *propertyRepository) Create(ctx context.Context, property *models.Property) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create property first
		if err := tx.Create(property).Error; err != nil {
			return fmt.Errorf("failed to create property: %w", err)
		}

		// If address is provided, ensure PropertyID is set and create it
		if property.Address != nil {
			property.Address.PropertyID = property.ID
			if err := tx.Create(property.Address).Error; err != nil {
				return fmt.Errorf("failed to create property address: %w", err)
			}
		}

		return nil
	})
}

// GetByID retrieves a Property by its unique ID, preloading related entities.
// ctx: context for the operation.
// id: UUID of the property.
// Returns: pointer to Property and error if not found or on failure.
func (r *propertyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Property, error) {
	var property models.Property
	err := r.db.
		Preload("Address").
		Preload("AssessmentDetails").
		Preload("Amenities").
		Preload("ConstructionDetails").
		Preload("ConstructionDetails.FloorDetails").
		Preload("AdditionalDetails").
		Preload("GISData").
		Preload("GISData.Coordinates").
		Preload("IGRS").
		Preload("Documents").
		Where("id = ?", id).First(&property).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("property with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get property by id %s: %w", id, err)
	}
	return &property, nil
}

// Update modifies an existing Property record in the database.
// ctx: context for the operation.
// property: pointer to Property model with updated data.
// Returns: error if update fails or record not found.
func (r *propertyRepository) Update(ctx context.Context, property *models.Property) error {
	result := r.db.Model(property).Updates(property)
    if result.Error != nil {
        return fmt.Errorf("failed to update property with id %s: %w", property.ID, result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("property with id %s not found for update", property.ID)
    }
    return nil
}

// Delete removes a Property record by its unique ID.
// ctx: context for the operation.
// id: UUID of the property to delete.
// Returns: error if deletion fails or record not found.
func (r *propertyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.Property{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete property with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("property with id %s not found for deletion", id)
	}
	return nil
}

// GetAll retrieves all Property records, optionally filtered by propertyType, with pagination and preloaded relations.
// ctx: context for the operation.
// page: page number (zero-based), size: number of records per page, propertyType: optional filter.
// Returns: slice of Property pointers, total count, and error if any.
func (r *propertyRepository) GetAll(ctx context.Context, page, size int, propertyType *string) ([]*models.Property, int64, error) {
	var properties []*models.Property
	var total int64

	query := r.db.Model(&models.Property{}).
		Preload("Address").
		Preload("AssessmentDetails").
		Preload("Amenities").
		Preload("ConstructionDetails").
		Preload("ConstructionDetails.FloorDetails").
		Preload("AdditionalDetails").
		Preload("GISData").
		Preload("GISData.Coordinates").
		Preload("IGRS").
		Preload("Documents")

	if propertyType != nil && *propertyType != "" {
		query = query.Where("property_type = ?", *propertyType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count properties: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Find(&properties).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no properties found")
		}
		return nil, 0, fmt.Errorf("failed to get properties: %w", err)
	}

	return properties, total, nil
}

// GetByPropertyNo retrieves a Property by its property number, preloading related entities.
// ctx: context for the operation.
// propertyNo: property number string.
// Returns: pointer to Property and error if not found or on failure.
func (r *propertyRepository) GetByPropertyNo(ctx context.Context, propertyNo string) (*models.Property, error) {
	var property models.Property
	err := r.db.
		Preload("Address").
		Preload("AssessmentDetails").
		Preload("Amenities").
		Preload("ConstructionDetails").
		Preload("ConstructionDetails.FloorDetails").
		Preload("AdditionalDetails").
		Preload("GISData").
		Preload("GISData.Coordinates").
		Preload("Documents").
		Preload("IGRS").
		Where("property_no = ?", propertyNo).First(&property).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("property with property_no %s not found", propertyNo)
		}
		return nil, fmt.Errorf("failed to get property by property_no %s: %w", propertyNo, err)
	}
	return &property, nil
}

// Search finds Property records matching the given search parameters, with filtering, sorting, pagination, and preloaded relations.
// ctx: context for the operation.
// params: SearchPropertyParams struct with filter and pagination options.
// Returns: slice of Property pointers, total count, and error if any.
func (r *propertyRepository) Search(ctx context.Context, params SearchPropertyParams) ([]*models.Property, int64, error) {
	var properties []*models.Property
	var total int64

	// Use the exact quoted table names as they appear in PostgreSQL
	propertiesTable := `"DIGIT3"."properties"`
	addressesTable := `"DIGIT3"."property_addresses"`

	query := r.db.Model(&models.Property{}).
		Preload("Address").
		Preload("AssessmentDetails").
		Preload("ConstructionDetails").
		Preload("AdditionalDetails").
		Joins(fmt.Sprintf("LEFT JOIN %s pa ON %s.id = pa.property_id", addressesTable, propertiesTable))

	// Apply filters
	var conditions []string
	var args []interface{}

	if params.PropertyType != nil && *params.PropertyType != "" {
		conditions = append(conditions, fmt.Sprintf("%s.property_type = ?", propertiesTable))
		args = append(args, *params.PropertyType)
	}

	if params.OwnershipType != nil && *params.OwnershipType != "" {
		conditions = append(conditions, fmt.Sprintf("%s.ownership_type = ?", propertiesTable))
		args = append(args, *params.OwnershipType)
	}

	if params.ComplexName != nil && *params.ComplexName != "" {
		conditions = append(conditions, fmt.Sprintf("%s.complex_name ILIKE ?", propertiesTable))
		args = append(args, "%"+*params.ComplexName+"%")
	}

	if params.Locality != nil && *params.Locality != "" {
		conditions = append(conditions, "pa.locality ILIKE ?")
		args = append(args, "%"+*params.Locality+"%")
	}

	if params.WardNo != nil && *params.WardNo != "" {
		conditions = append(conditions, "pa.ward_no = ?")
		args = append(args, *params.WardNo)
	}

	if params.ZoneNo != nil && *params.ZoneNo != "" {
		conditions = append(conditions, "pa.zone_no = ?")
		args = append(args, *params.ZoneNo)
	}

	if params.Street != nil && *params.Street != "" {
		conditions = append(conditions, "pa.street ILIKE ?")
		args = append(args, "%"+*params.Street+"%")
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count properties: %w", err)
	}

	// Apply sorting with proper table names
	orderBy := fmt.Sprintf("%s.created_at DESC", propertiesTable)
	if params.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(params.SortOrder) == "DESC" {
			direction = "DESC"
		}

		switch params.SortBy {
		case "createdAt":
			orderBy = fmt.Sprintf("%s.created_at %s", propertiesTable, direction)
		case "updatedAt":
			orderBy = fmt.Sprintf("%s.updated_at %s", propertiesTable, direction)
		case "propertyNo":
			orderBy = fmt.Sprintf("%s.property_no %s", propertiesTable, direction)
		case "propertyType":
			orderBy = fmt.Sprintf("%s.property_type %s", propertiesTable, direction)
		default:
			orderBy = fmt.Sprintf("%s.created_at DESC", propertiesTable)
		}
	}

	// Apply pagination and ordering
	offset := params.Page * params.Size
	if err := query.Order(orderBy).Offset(offset).Limit(params.Size).Find(&properties).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no properties found for search params")
		}
		return nil, 0, fmt.Errorf("failed to search properties: %w", err)
	}

	return properties, total, nil
}
