package handlers

import (
	"bytes"
	"encoding/json"
	"enumeration/internal/models"
	"enumeration/internal/services"
	"enumeration/pkg/response"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DocumentHandler handles HTTP requests for document resources.
type DocumentHandler struct {
	service services.DocumentService // Service layer for document operations
}

// NewDocumentHandler creates a new DocumentHandler with the provided service.
func NewDocumentHandler(service services.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

// Create handles POST requests to create a single document record.
// Validates the request body and delegates creation to the service layer.
func (h *DocumentHandler) Create(c *gin.Context) {
	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.CreateDocument(c.Request.Context(), &doc); err != nil {
		// Map validation errors to 400, others to 500
		if errors.Is(err, services.ErrValidation) {
			response.BadRequest(c, "Invalid document: "+err.Error())
			return
		}
		response.InternalServerError(c, "Failed to create document: "+err.Error())
		return
	}

	response.Created(c, response.SuccessResponseBody("Document created successfully", doc))
}

// CreateBatch handles POST requests to create multiple documents (batch).
// Accepts either a single Document object or an array of Document objects in the request body.
func (h *DocumentHandler) CreateBatch(c *gin.Context) {
	// Read raw body so we can unmarshal conditionally
	raw, err := c.GetRawData()
	if err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Trim leading whitespace to inspect first non-space char
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	var docs []*models.Document

	if len(trimmed) == 0 {
		response.BadRequest(c, "Empty request body")
		return
	}

	// If JSON array -> unmarshal into slice
	if trimmed[0] == '[' {
		if err := json.Unmarshal(raw, &docs); err != nil {
			response.BadRequest(c, "Invalid request body: "+err.Error())
			return
		}
	} else { // JSON object -> unmarshal single doc and wrap into slice
		var doc models.Document
		if err := json.Unmarshal(raw, &doc); err != nil {
			response.BadRequest(c, "Invalid request body: "+err.Error())
			return
		}
		docs = append(docs, &doc)
	}

	if err := h.service.CreateDocuments(c.Request.Context(), docs); err != nil {
		if errors.Is(err, services.ErrValidation) {
			response.BadRequest(c, "Invalid documents: "+err.Error())
			return
		}
		response.InternalServerError(c, "Failed to create documents: "+err.Error())
		return
	}

	response.Created(c, response.SuccessResponseBody("Documents created successfully", docs))
}

// GetByID handles GET requests to retrieve a document by its ID.
func (h *DocumentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid document ID")
		return
	}

	doc, err := h.service.GetDocumentByID(c.Request.Context(), id)
	if err != nil {
		// validation error -> 400, otherwise map not found to 404 if service/repo returns that
		if errors.Is(err, services.ErrValidation) {
			response.BadRequest(c, "Invalid document ID: "+err.Error())
			return
		}
		response.NotFound(c, "Document not found")
		return
	}

	response.Success(c, "Document retrieved successfully", doc)
}

// GetByPropertyID handles GET requests to retrieve all documents for a given property ID.
func (h *DocumentHandler) GetByPropertyID(c *gin.Context) {
	propertyIDStr := c.Param("propertyId")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid property ID")
		return
	}

	docs, err := h.service.GetDocumentsByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			response.BadRequest(c, "Invalid property ID: "+err.Error())
			return
		}
		response.InternalServerError(c, "Failed to retrieve documents: "+err.Error())
		return
	}

	response.Success(c, "Documents retrieved successfully", docs)
}

// GetAll handles GET requests to retrieve all documents with pagination.
// Optionally filters by property ID and sets pagination headers.
func (h *DocumentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var propertyID *uuid.UUID
	if propertyIDStr := c.Query("propertyId"); propertyIDStr != "" {
		if id, err := uuid.Parse(propertyIDStr); err == nil {
			propertyID = &id
		} else {
			response.BadRequest(c, "Invalid propertyId query parameter")
			return
		}
	}

	docs, total, err := h.service.GetAllDocuments(c.Request.Context(), page, size, propertyID)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve documents: "+err.Error())
		return
	}

	response.Paginated(c, docs, response.PaginationMeta{
		Page:       page,
		PageSize:   size,
		TotalItems: total,
		TotalPages: int((total + int64(size) - 1) / int64(size)),
	})
}

// Delete handles DELETE requests to remove a document by its ID.
func (h *DocumentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid document ID")
		return
	}

	if err := h.service.DeleteDocument(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrValidation) {
			response.BadRequest(c, "Invalid document ID: "+err.Error())
			return
		}
		response.InternalServerError(c, "Failed to delete document: "+err.Error())
		return
	}

	response.Success(c, "Document deleted successfully", nil)
}

// Update handles PUT requests to update an existing document by its ID.
// Validates the request body and delegates update to the service layer.
func (h *DocumentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid document ID")
		return
	}

	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		response.BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	doc.ID = id

	if err := h.service.UpdateDocument(c.Request.Context(), &doc); err != nil {
		if errors.Is(err, services.ErrValidation) {
			response.BadRequest(c, "Invalid document: "+err.Error())
			return
		}
		response.InternalServerError(c, "Failed to update document: "+err.Error())
		return
	}

	response.Success(c, "Document updated successfully", doc)
}
