package repositories

import (
	"context"
	"enumeration/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// documentRepository implements the DocumentRepository interface using GORM.
var _ DocumentRepository = (*documentRepository)(nil)

type documentRepository struct {
	db *gorm.DB // GORM database connection
}

// NewDocumentRepository creates a new repository for documents.
func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

// Create inserts a new Document record into the database.
func (r *documentRepository) Create(ctx context.Context, doc *models.Document) error {
	if err := r.db.Create(doc).Error; err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	return nil
}

// CreateBatch inserts multiple Document records in a single transaction for efficiency and rollback on failure.
func (r *documentRepository) CreateBatch(ctx context.Context, docs []*models.Document) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, d := range docs {
			if err := tx.Create(d).Error; err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}
		}
		return nil
	})
}

// GetByID retrieves a Document record by its ID.
func (r *documentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	var doc models.Document
	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch document by id %s: %w", id, err)
	}
	return &doc, nil
}

// GetByPropertyID retrieves all Document records for a given property ID.
func (r *documentRepository) GetByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.Document, error) {
	var docs []*models.Document
	err := r.db.Where("property_id = ?", propertyID).Find(&docs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no documents found for property id %s", propertyID)
		}
		return nil, fmt.Errorf("failed to get documents by property id %s: %w", propertyID, err)
	}
	return docs, nil
}

// GetAll retrieves all Document records with optional filtering by property ID and supports pagination.
// Returns the filtered documents and the total count.
func (r *documentRepository) GetAll(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.Document, int64, error) {
	var docs []*models.Document
	var total int64

	query := r.db.Model(&models.Document{})

	if propertyID != nil {
		query = query.Where("property_id = ?", *propertyID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	offset := page * size
	err := query.Offset(offset).Limit(size).Find(&docs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("no documents found")
		}
		return nil, 0, fmt.Errorf("failed to get documents: %w", err)
	}

	return docs, total, nil
}

// Delete removes a Document record by its ID.
// Returns an error if the record does not exist or deletion fails.
func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&models.Document{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete document with id %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("document with id %s not found for deletion", id)
	}
	return nil
}

// Update modifies an existing Document record by its ID.
// Returns an error if the record does not exist or update fails.
func (r *documentRepository) Update(ctx context.Context, doc *models.Document) error {
	result := r.db.Model(&models.Document{}).Where("id = ?", doc.ID).Updates(doc)
	if result.Error != nil {
		return fmt.Errorf("failed to update document with id %s: %w", doc.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("document with id %s not found for update", doc.ID)
	}
	return nil
}
