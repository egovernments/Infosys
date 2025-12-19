package repositories

import (
	"context"
	"enumeration/internal/dto"
	"enumeration/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ApplicationRepository = (*applicationRepository)(nil)

// applicationRepository handles database operations for applications
type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository creates a new application repository
func NewApplicationRepository(db *gorm.DB) *applicationRepository {
	return &applicationRepository{db: db}
}

// Create creates a new application record
func (r *applicationRepository) Create(ctx context.Context, application *models.Application) error {
	if err := r.db.Create(application).Error; err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	return nil
}

// GetByID retrieves application by ID
func (r *applicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	var application models.Application
	err := r.db.Preload("Property").
		Preload("Property.Address").
		Preload("Property.AssessmentDetails").
		Preload("Property.Amenities").
		Preload("Property.ConstructionDetails").
		Preload("Property.ConstructionDetails.FloorDetails").
		Preload("Property.AdditionalDetails").
		Preload("Property.GISData").
		Preload("Property.GISData.Coordinates").
		Preload("Property.IGRS").
		Preload("ApplicationLogs").
		Preload("Property.Documents").
		First(&application, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get application by id %s: %w", id, err)
	}
	return &application, nil
}

// GetByAssignedAgent retrieves applications by assigned agent ID with pagination
func (r *applicationRepository) GetByAssignedAgent(ctx context.Context, agentID string, status string, page, size int) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	// Build query with filters
	query := r.db.WithContext(ctx).Model(&models.Application{}).
		Preload("Property").
		Preload("Property.Address").
		Preload("Property.AssessmentDetails").
		Preload("Property.Amenities").
		Preload("Property.ConstructionDetails").
		Preload("Property.ConstructionDetails.FloorDetails").
		Preload("Property.AdditionalDetails").
		Preload("Property.GISData").
		Preload("Property.GISData.Coordinates").
		Preload("ApplicationLogs").
		Preload("Property.Documents")

	// Add assigned agent filter
	query = query.Where("assigned_agent = ?", agentID)

	query = query.Where("status = ?", status)
	query = query.Where("is_draft = ?", false)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count applications for agent %s: %w", agentID, err)
	}

	// Apply pagination and ordering
	offset := page * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&applications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("no applications found for agent %s with status %s", agentID, status)
		}
		return nil, 0, fmt.Errorf("failed to get applications for agent %s: %w", agentID, err)
	}
	return applications, total, nil
}

// GetByTenantIDAndStatus retrieves applications by tenant ID and status with pagination
func (r *applicationRepository) GetByTenantIDAndStatus(ctx context.Context, tenantID string, status string, page, size int) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	// Build query with filters
	query := r.db.WithContext(ctx).Model(&models.Application{}).
		Preload("Property").
		Preload("Property.Address").
		Preload("Property.AssessmentDetails").
		Preload("Property.Amenities").
		Preload("Property.ConstructionDetails").
		Preload("Property.ConstructionDetails.FloorDetails").
		Preload("Property.AdditionalDetails").
		Preload("Property.GISData").
		Preload("Property.GISData.Coordinates").
		Preload("Property.IGRS").
		Preload("ApplicationLogs").
		Preload("Property.Documents")

	// Add tenant ID filter
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// Add status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query = query.Where("is_draft = ?", false)
	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count applications for tenant %s: %w", tenantID, err)
	}

	// Apply pagination and ordering
	offset := page * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get applications for tenant %s: %w", tenantID, err)
	}
	return applications, total, nil
}

// GetByApplicationNo retrieves application by application number
func (r *applicationRepository) GetByApplicationNo(ctx context.Context, applicationNo string) (*models.Application, error) {
	var application models.Application
	res := r.db.First(&application, "application_no = ?", applicationNo)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to get application by application number %s: %w", applicationNo, res.Error)
	}
	return &application, nil
}

// Update updates an existing application record
func (r *applicationRepository) Update(ctx context.Context, application *models.Application) error {
	result := r.db.Save(application)
	if result.Error != nil {
		return fmt.Errorf("failed to update application with id %s: %w", application.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application with id %s not found for update", application.ID)
	}
	return nil
}

// Delete deletes application by ID
func (r *applicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.Application{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete application with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application with id %s not found for deletion", id)
	}
	return nil
}

// List retrieves all applications with pagination
func (r *applicationRepository) List(ctx context.Context, page, size int) ([]models.Application, int64, error) {
	var applications []models.Application
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Application{}).
		Preload("Property").
		Preload("Property.Address").
		Preload("Property.AssessmentDetails").
		Preload("Property.Amenities").
		Preload("Property.ConstructionDetails").
		Preload("Property.ConstructionDetails.FloorDetails").
		Preload("Property.AdditionalDetails").
		Preload("Property.GISData").
		Preload("Property.GISData.Coordinates").
		Preload("Property.IGRS").
		Preload("Property.Documents")
	query = query.Where("is_draft = ?", false)
	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}

	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
	}

	return applications, total, nil
}

// Search retrieves applications based on search criteria
// func (r *applicationRepository) Search(ctx context.Context, criteria *dto.ApplicationSearchCriteria, page, size int) ([]*models.Application, int64, error) {
// 	var applications []*models.Application
// 	var total int64

// 	query := r.db.WithContext(ctx).Model(&models.Application{}).Preload("Property").
// 		Preload("Property.Address").
// 		Preload("Property.AssessmentDetails").
// 		Preload("Property.Amenities").
// 		Preload("Property.ConstructionDetails").
// 		Preload("Property.ConstructionDetails.FloorDetails").
// 		Preload("Property.AdditionalDetails").
// 		Preload("Property.GISData").
// 		Preload("Property.GISData.Coordinates").
// 		Preload("Property.IGRS").
// 		Preload("Property.Documents")
// 	// Apply filters
// 	if criteria.Status != "" {
// 		query = query.Where("status = ?", criteria.Status)
// 	}
// 	if criteria.Priority != "" {
// 		query = query.Where("priority = ?", criteria.Priority)
// 	}
// 	if criteria.PropertyID != "" {
// 		query = query.Where("property_id = ?", criteria.PropertyID)
// 	}
// 	if criteria.AssignedAgent != "" {
// 		query = query.Where("assigned_agent = ?", criteria.AssignedAgent)
// 	}
// 	if criteria.AppliedBy != "" {
// 		query = query.Where("applied_by = ?", criteria.AppliedBy)
// 	}
// 	if criteria.ApplicationNo != "" {
// 		query = query.Where("application_no = ?", criteria.ApplicationNo)
// 	}
// 	if criteria.CreatedDateFrom != nil {
// 		query = query.Where("created_at >= ?", *criteria.CreatedDateFrom)
// 	}
// 	if criteria.CreatedDateTo != nil {
// 		query = query.Where("created_at <= ?", *criteria.CreatedDateTo)
// 	}
// 	if criteria.DueDateFrom != nil {
// 		query = query.Where("due_date >= ?", *criteria.DueDateFrom)
// 	}
// 	if criteria.DueDateTo != nil {
// 		query = query.Where("due_date <= ?", *criteria.DueDateTo)
// 	}
// 	if criteria.AssesseeID != "" {
// 		query = query.Where("assessee_id = ?", criteria.AssesseeID)
// 	}
// 	if criteria.ZoneNo != "" && criteria.WardNo != "" {
// 		query = query.Joins(`JOIN "DIGIT3"."properties" ON "DIGIT3"."applications"."property_id" = "DIGIT3"."properties"."id"`).
// 			Joins(`JOIN "DIGIT3"."property_addresses" ON "DIGIT3"."properties"."id" = "DIGIT3"."property_addresses"."property_id"`).
// 			Where(`"DIGIT3"."property_addresses"."zone_no" = ? AND "DIGIT3"."property_addresses"."ward_no" = ?`, criteria.ZoneNo, criteria.WardNo)
// 	} else if criteria.ZoneNo != "" {
// 		query = query.Joins(`JOIN "DIGIT3"."properties" ON "DIGIT3"."applications"."property_id" = "DIGIT3"."properties"."id"`).
// 			Joins(`JOIN "DIGIT3"."property_addresses" ON "DIGIT3"."properties"."id" = "DIGIT3"."property_addresses"."property_id"`).
// 			Where(`"DIGIT3"."property_addresses"."zone_no" = ?`, criteria.ZoneNo)
// 	} else if criteria.WardNo != "" {
// 		query = query.Joins(`JOIN "DIGIT3"."properties" ON "DIGIT3"."applications"."property_id" = "DIGIT3"."properties"."id"`).
// 			Joins(`JOIN "DIGIT3"."property_addresses" ON "DIGIT3"."properties"."id" = "DIGIT3"."property_addresses"."property_id"`).
// 			Where(`"DIGIT3"."property_addresses"."ward_no" = ?`, criteria.WardNo)
// 	}
// 	if criteria.IsDraft != nil {
// 		query = query.Where("is_draft = ?", *criteria.IsDraft)
// 	} else {
// 		query = query.Where("is_draft = ?", false)
// 	}
// 	// Count total records
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
// 	}

// 	// Apply pagination and sorting
// 	offset := page * size

// 	// Determine sort field (default to created_at)
// 	sortField := "created_at"
// 	if criteria.SortField != "" {
// 		sortField = criteria.SortField
// 	}

// 	// Determine sort order (default to DESC)
// 	sortOrder := "DESC"
// 	if criteria.SortBy == "ASC" {
// 		sortOrder = "ASC"
// 	}

// 	// Build order clause
// 	orderClause := fmt.Sprintf("%s %s", sortField, sortOrder)

// 	if err := query.Offset(offset).Limit(size).Order(orderClause).Find(&applications).Error; err != nil {
// 		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
// 	}
// 	return applications, total, nil
// }

// Search retrieves applications based on search criteria
func (r *applicationRepository) Search(ctx context.Context, criteria *dto.ApplicationSearchCriteria, page, size int) ([]*models.Application, int64, error) {
    var applications []*models.Application
    var total int64

    query := r.db.WithContext(ctx).Model(&models.Application{}).Preload("Property").
        Preload("Property.Address").
        Preload("Property.AssessmentDetails").
        Preload("Property.Amenities").
        Preload("Property.ConstructionDetails").
        Preload("Property.ConstructionDetails.FloorDetails").
        Preload("Property.AdditionalDetails").
        Preload("Property.GISData").
        Preload("Property.GISData.Coordinates").
        Preload("Property.IGRS").
        Preload("Property.Documents")
        
    // Apply filters
    if criteria.Status != "" {
        query = query.Where("status = ?", criteria.Status)
    }
    if criteria.Priority != "" {
        query = query.Where("priority = ?", criteria.Priority)
    }
    if criteria.PropertyID != "" {
        query = query.Where("property_id = ?", criteria.PropertyID)
    }
    if criteria.AssignedAgent != "" {
        query = query.Where("assigned_agent = ?", criteria.AssignedAgent)
    }
    if criteria.AppliedBy != "" {
        query = query.Where("applied_by = ?", criteria.AppliedBy)
    }
    if criteria.ApplicationNo != "" {
        query = query.Where("application_no = ?", criteria.ApplicationNo)
    }
    if criteria.CreatedDateFrom != nil {
        query = query.Where("created_at >= ?", *criteria.CreatedDateFrom)
    }
    if criteria.CreatedDateTo != nil {
        query = query.Where("created_at <= ?", *criteria.CreatedDateTo)
    }
    if criteria.DueDateFrom != nil {
        query = query.Where("due_date >= ?", *criteria.DueDateFrom)
    }
    if criteria.DueDateTo != nil {
        query = query.Where("due_date <= ?", *criteria.DueDateTo)
    }
    if criteria.AssesseeID != "" {
        query = query.Where("assessee_id = ?", criteria.AssesseeID)
    }

    // Updated logic for PropertyNo, ZoneNo and WardNo
    needsJoin := false
    
    if criteria.PropertyNo != "" {
        needsJoin = true
        fmt.Println("PropertyNo filter detected, will add joins")
    }
    if criteria.ZoneNo != "" {
        needsJoin = true
    }
    if len(criteria.WardNo) > 0 {
        needsJoin = true
    }

    // Add joins only if needed
    if needsJoin {
        query = query.Joins(`JOIN "DIGIT3"."properties" ON "DIGIT3"."applications"."property_id" = "DIGIT3"."properties"."id"`)
        fmt.Println("Added properties join")
        
        // Only join address table if zone or ward filters are present
        if criteria.ZoneNo != "" || len(criteria.WardNo) > 0 {
            query = query.Joins(`JOIN "DIGIT3"."property_addresses" ON "DIGIT3"."properties"."id" = "DIGIT3"."property_addresses"."property_id"`)
        }
    }

    // Apply property number filter
    if criteria.PropertyNo != "" {
        query = query.Where(`"DIGIT3"."properties"."property_no" = ?`, criteria.PropertyNo)
        fmt.Printf("Added PropertyNo filter: %s\n", criteria.PropertyNo)
    }

    // Apply zone filter
    if criteria.ZoneNo != "" {
        query = query.Where(`"DIGIT3"."property_addresses"."zone_no" = ?`, criteria.ZoneNo)
    }

    // Apply ward filter - handle array of ward numbers
    if len(criteria.WardNo) > 0 {
        // Filter out empty strings from the slice
        var validWards []string
        for _, ward := range criteria.WardNo {
            if ward != "" {
                validWards = append(validWards, ward)
            }
        }
        
        if len(validWards) > 0 {
            query = query.Where(`"DIGIT3"."property_addresses"."ward_no" IN ?`, validWards)
        }
    }

    if criteria.IsDraft != nil {
        query = query.Where("is_draft = ?", *criteria.IsDraft)
    } else {
        query = query.Where("is_draft = ?", false)
    }

    // Debug: Print the SQL query
    sqlQuery := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
        return tx.Count(&total)
    })
    fmt.Printf("Count SQL: %s\n", sqlQuery)

    // Count total records
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count applications: %w", err)
    }

    fmt.Printf("Total records found: %d\n", total)

    // Apply pagination and sorting
    offset := page * size

    // Determine sort field (default to created_at)
    sortField := "created_at"
    if criteria.SortField != "" {
        sortField = criteria.SortField
    }

    // Determine sort order (default to DESC)
    sortOrder := "DESC"
    if criteria.SortBy == "ASC" {
        sortOrder = "ASC"
    }

    // Build order clause - ensure proper table qualification for joins
    var orderClause string
    if needsJoin && (sortField == "created_at" || sortField == "due_date") {
        // Qualify the sort field with table name when joins are present
        orderClause = fmt.Sprintf(`"DIGIT3"."applications"."%s" %s`, sortField, sortOrder)
    } else {
        orderClause = fmt.Sprintf("%s %s", sortField, sortOrder)
    }

    if err := query.Offset(offset).Limit(size).Order(orderClause).Find(&applications).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to get applications: %w", err)
    }
    return applications, total, nil
}

// ExistsByApplicationNo checks if application exists by application number
func (r *applicationRepository) ExistsByApplicationNo(ctx context.Context, applicationNo string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Application{}).Where("application_no = ?", applicationNo).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check existence of application number %s: %w", applicationNo, err)
	}
	return count > 0, nil
}

func (r *applicationRepository) GetApplicationsByIDs(ctx context.Context, ids []string, state string, page, size int) ([]models.Application, int64, error) {
	if len(ids) == 0 {
		return []models.Application{}, 0, nil
	}
	var apps []models.Application
	// Build query with filters
	query := r.db.WithContext(ctx).Preload("Property").
		Preload("Property.Address").
		Preload("Property.AssessmentDetails").
		Preload("Property.Amenities").
		Preload("Property.ConstructionDetails").
		Preload("Property.ConstructionDetails.FloorDetails").
		Preload("Property.AdditionalDetails").
		Preload("Property.Documents").
		Preload("Property.GISData").
		Preload("Property.GISData.Coordinates").
		Preload("Property.GISData").Where("id IN ?", ids)

	query = query.Where("status = ?", state)

	// Get total count
	var totalCount int64
	if err := query.Model(&models.Application{}).Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}
	// Apply pagination
	offset := page * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&apps).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
	}

	return apps, totalCount, nil
}

func (r *applicationRepository) CheckUserExists(ctx context.Context, userId string) error {
	var count int64
	err := r.db.Table("DIGIT3.users").Where("keycloak_user_id = ?", userId).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("user with ID %s does not exist", userId)
	}
	return nil
}
