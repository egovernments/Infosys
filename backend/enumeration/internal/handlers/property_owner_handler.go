package handlers

import (
	"bytes"
	"encoding/json"
	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PropertyOwnerHandler handles HTTP requests for property owner resources.
type PropertyOwnerHandler struct {
	service services.PropertyOwnerService // Service layer for property owner operations
}

// NewPropertyOwnerHandler creates a new PropertyOwnerHandler with the provided service.
func NewPropertyOwnerHandler(service services.PropertyOwnerService) *PropertyOwnerHandler {
	return &PropertyOwnerHandler{service: service}
}

// Create handles POST /property-owners (single).
// Validates the request body and delegates creation to the service layer.
func (h *PropertyOwnerHandler) Create(c *gin.Context) {
	var req dto.CreatePropertyOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
		return
	}

	owner, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid property owner", "errors": []string{err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create property owner", "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Property owner created successfully", "data": owner})
}

// CreateBatch handles POST /property-owners/batch and accepts either a single object or an array.
// Supports both single and batch creation of property owners by detecting the request body type.
func (h *PropertyOwnerHandler) CreateBatch(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
		return
	}
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	if len(trimmed) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Empty request body"})
		return
	}

	var reqPtrs []*dto.CreatePropertyOwnerRequest
	if trimmed[0] == '[' {
		// Handle array of property owners
		var arr []dto.CreatePropertyOwnerRequest
		if err := json.Unmarshal(raw, &arr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
			return
		}
		reqPtrs = make([]*dto.CreatePropertyOwnerRequest, 0, len(arr))
		for i := range arr {
			reqPtrs = append(reqPtrs, &arr[i])
		}
	} else {
		// Handle single property owner
		var single dto.CreatePropertyOwnerRequest
		if err := json.Unmarshal(raw, &single); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
			return
		}
		reqPtrs = []*dto.CreatePropertyOwnerRequest{&single}
	}

	owners, err := h.service.CreatePropertyOwners(c.Request.Context(), reqPtrs)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid property owners payload", "errors": []string{err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create property owners", "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Property owners created successfully", "data": owners})
}

// GetByPropertyID handles GET /property-owners/property/:propertyId?page=0&size=20
// Retrieves property owners for a given property ID with pagination support.
// Sets pagination headers in the response.
func (h *PropertyOwnerHandler) GetByPropertyID(c *gin.Context) {
	propertyID, err := uuid.Parse(c.Param("propertyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid property ID", "errors": []string{err.Error()}})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    if size <= 0 {
		size = 20
	}
	owners, total, err := h.service.GetByPropertyID(c.Request.Context(), propertyID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to retrieve property owners", "errors": []string{err.Error()}})
		return
	}

	// Calculate total pages for pagination
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}

	c.Header(constants.HeaderTotalCount, strconv.FormatInt(total, 10))
	c.Header(constants.HeaderCurrentPage, strconv.Itoa(page))
	c.Header(constants.HeaderPerPage, strconv.Itoa(size))
	c.Header(constants.HeaderTotalPages, strconv.Itoa(totalPages))

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Property owners retrieved successfully", "data": owners})
}

// Update handles PUT /property-owners/:id
// Updates an existing property owner by ID. Validates the request body and delegates update to the service layer.
func (h *PropertyOwnerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid property owner ID", "errors": []string{err.Error()}})
		return
	}

	var req dto.UpdatePropertyOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
		return
	}

	owner, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid update payload", "errors": []string{err.Error()}})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to update property owner", "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Property owner updated successfully", "data": owner})
}

// Delete handles DELETE /property-owners/:id
// Deletes a property owner by ID. Delegates deletion to the service layer.
func (h *PropertyOwnerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid property owner ID", "errors": []string{err.Error()}})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Failed to delete property owner", "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Property owner deleted successfully"})
}
