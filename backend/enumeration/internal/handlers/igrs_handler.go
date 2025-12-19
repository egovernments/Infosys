// Package handlers contains HTTP handler implementations for the property tax enumeration system.
package handlers

import (
	"enumeration/internal/dto"
	"enumeration/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IGRSHandler handles HTTP requests for IGRS (Integrated Grievance Redressal System) resources.
type IGRSHandler struct {
	service services.IGRSService // Service layer for IGRS operations
}

// NewIGRSHandler creates a new IGRSHandler with the provided service.
func NewIGRSHandler(s services.IGRSService) *IGRSHandler {
	return &IGRSHandler{service: s}
}

// Create handles POST requests to create a new IGRS record.
// Validates the request body and delegates creation to the service layer.
func (h *IGRSHandler) Create(c *gin.Context) {
	var req dto.CreateIGRSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
		return
	}
	out, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create", "errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "IGRS created", "data": out})
}

// GetByID handles GET requests to retrieve an IGRS record by its ID.
func (h *IGRSHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid id", "errors": []string{err.Error()}})
		return
	}
	out, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Not found", "errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": out})
}

// Update handles PUT requests to update an existing IGRS record by its ID.
// Validates the request body and delegates update to the service layer.
func (h *IGRSHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid id", "errors": []string{err.Error()}})
		return
	}
	var req dto.UpdateIGRSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "errors": []string{err.Error()}})
		return
	}
	out, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update", "errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Updated", "data": out})
}

// Delete handles DELETE requests to remove an IGRS record by its ID.
func (h *IGRSHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid id", "errors": []string{err.Error()}})
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete", "errors": []string{err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted"})
}

// List handles GET requests to retrieve a paginated list of IGRS records.
// Sets pagination headers and returns the result set.
func (h *IGRSHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	items, total, err := h.service.List(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to list", "errors": []string{err.Error()}})
		return
	}
	totalPages := int(total) / size
	if int(total)%size != 0 {
		totalPages++
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}
