package services

import (
	"context"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Compile-time check for interface implementation
var _ DocumentService = (*documentService)(nil)
var ErrValidation = errors.New("validation error")

// documentService handles business logic for documents
type documentService struct {
	repo repositories.DocumentRepository
}

// NewDocumentService returns a new documentService
func NewDocumentService(repo repositories.DocumentRepository) DocumentService {
	return &documentService{repo: repo}
}

// CreateDocument validates and adds a new document
func (s *documentService) CreateDocument(ctx context.Context, doc *models.Document) error {
	// basic validation
	if doc == nil {
		return fmt.Errorf("%w: document is nil", ErrValidation)
	}

	if doc.PropertyID == uuid.Nil {
		return fmt.Errorf("%w: property ID is required", ErrValidation)
	}
	if doc.DocumentType == "" {
		return fmt.Errorf("%w: document type is required", ErrValidation)
	}
	if doc.DocumentName == "" {
		return fmt.Errorf("%w: document name is required", ErrValidation)
	}
	if doc.Action == "" {
		return fmt.Errorf("%w: action is required", ErrValidation)
	}

	// set upload date if not provided
	if doc.UploadDate.IsZero() {
		doc.UploadDate = time.Now().UTC()
	}

	// populate UploadedBy from context if available
	if doc.UploadedBy == "" {
		var username, role string

		// Extract username using the same context key type
		if v := ctx.Value("user"); v != nil {
			username = fmt.Sprint(v)
		}

		// Extract role using the same context key type
		if v := ctx.Value("role"); v != nil {
			role = fmt.Sprint(v)
		}

		// Format as "username role:rolename"
		if username != "" && role != "" {
			doc.UploadedBy = fmt.Sprintf("%s role:%s", username, role)
		} else if username != "" {
			doc.UploadedBy = username
		}
	}

	// ensure ID is not nil (DB default will set it, but keep consistency)
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}

	return s.repo.Create(ctx, doc)
}

// CreateDocuments validates and adds multiple documents in a batch
func (s *documentService) CreateDocuments(ctx context.Context, docs []*models.Document) error {
	if docs == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	for i, d := range docs {
		if d == nil {
			return fmt.Errorf("%w: document at index %d is nil", ErrValidation, i)
		}
	}

	// try to extract username and role once from context to reuse
	var uploadedByFromCtx string
	var username, role string

	// Extract username using the same context key type
	if v := ctx.Value("user"); v != nil {
		username = fmt.Sprint(v)
	}

	// Extract role using the same context key type
	if v := ctx.Value("role"); v != nil {
		role = fmt.Sprint(v)
	}

	// Format as "username role:rolename"
	if username != "" && role != "" {
		uploadedByFromCtx = fmt.Sprintf("%s role:%s", username, role)
	} else if username != "" {
		uploadedByFromCtx = username
	}

	for _, d := range docs {
		if d.PropertyID == uuid.Nil {
			return fmt.Errorf("%w: property ID is required for all documents", ErrValidation)
		}
		if d.DocumentType == "" {
			return fmt.Errorf("%w: document type is required for all documents", ErrValidation)
		}
		if d.DocumentName == "" {
			return fmt.Errorf("%w: document name is required for all documents", ErrValidation)
		}
		if d.UploadDate.IsZero() {
			d.UploadDate = time.Now().UTC()
		}
		if d.ID == uuid.Nil {
			d.ID = uuid.New()
		}

		// populate UploadedBy if not already set on the document
		if d.UploadedBy == "" && uploadedByFromCtx != "" {
			d.UploadedBy = uploadedByFromCtx
		}
	}
	return s.repo.CreateBatch(ctx, docs)
}

// GetDocumentByID fetches a document by its ID
func (s *documentService) GetDocumentByID(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid document ID", ErrValidation)
	}
	return s.repo.GetByID(ctx, id)
}

// GetDocumentsByPropertyID fetches all documents for a property
func (s *documentService) GetDocumentsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]*models.Document, error) {
	if propertyID == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid property ID", ErrValidation)
	}
	return s.repo.GetByPropertyID(ctx, propertyID)
}

// GetAllDocuments returns all documents with pagination, optionally filtered by propertyID
func (s *documentService) GetAllDocuments(ctx context.Context, page, size int, propertyID *uuid.UUID) ([]*models.Document, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 || size > 100 {
		size = constants.DefaultSize
	}
	return s.repo.GetAll(ctx, page, size, propertyID)
}

// DeleteDocument removes a document by its ID
func (s *documentService) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid document ID", ErrValidation)
	}
	return s.repo.Delete(ctx, id)
}

// UpdateDocument updates an existing document after validation
func (s *documentService) UpdateDocument(ctx context.Context, doc *models.Document) error {
	if doc == nil {
		return fmt.Errorf("%w: document is nil", ErrValidation)
	}
	if doc.ID == uuid.Nil {
		return fmt.Errorf("%w: document ID is required", ErrValidation)
	}
	// Optionally, add more validation as needed
	return s.repo.Update(ctx, doc)
}
