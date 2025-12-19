// Package handlers contains HTTP handler implementations for the property tax enumeration system.
package handlers

import (
	"enumeration/internal/models"
	"enumeration/internal/services"
	"enumeration/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConstructionDetailsHandler handles HTTP requests for construction details resources.
type ConstructionDetailsHandler struct {
	constructionDetailsService services.ConstructionDetailsService // Service layer for construction details operations
}

// NewConstructionDetailsHandler creates a new ConstructionDetailsHandler with the provided service.
func NewConstructionDetailsHandler(constructionDetailsService services.ConstructionDetailsService) *ConstructionDetailsHandler {
	return &ConstructionDetailsHandler{
		constructionDetailsService: constructionDetailsService,
	}
}

// CreateConstructionDetails handles POST requests to create a new construction details record.
// Validates the request body and delegates creation to the service layer.
func (h *ConstructionDetailsHandler) CreateConstructionDetails(c *gin.Context) {
	var constructionDetails models.ConstructionDetails

	if err := c.ShouldBindJSON(&constructionDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	if err := h.constructionDetailsService.CreateConstructionDetails(c.Request.Context(), &constructionDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to create construction details", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponseBody("Construction details created successfully", constructionDetails))
}

// GetConstructionDetailsByID handles GET requests to retrieve construction details by their ID.
func (h *ConstructionDetailsHandler) GetConstructionDetailsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid construction details ID", err.Error()))
		return
	}

	constructionDetails, err := h.constructionDetailsService.GetConstructionDetailsByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "construction details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Construction details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get construction details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Construction details retrieved successfully", constructionDetails))
}

// UpdateConstructionDetails handles PUT requests to update existing construction details by their ID.
// Validates the request body and delegates update to the service layer.
func (h *ConstructionDetailsHandler) UpdateConstructionDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid construction details ID", err.Error()))
		return
	}

	var constructionDetails models.ConstructionDetails
	if err := c.ShouldBindJSON(&constructionDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	constructionDetails.ID = id

	if err := h.constructionDetailsService.UpdateConstructionDetails(c.Request.Context(), &constructionDetails); err != nil {
		if err.Error() == "construction details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Construction details not found", err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to update construction details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Construction details updated successfully", constructionDetails))
}

// DeleteConstructionDetails handles DELETE requests to remove construction details by their ID.
func (h *ConstructionDetailsHandler) DeleteConstructionDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid construction details ID", err.Error()))
		return
	}

	if err := h.constructionDetailsService.DeleteConstructionDetails(c.Request.Context(), id); err != nil {
		if err.Error() == "construction details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Construction details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to delete construction details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Construction details deleted successfully", nil))
}

// GetAllConstructionDetails handles GET requests to retrieve all construction details with pagination.
// Optionally filters by property ID and sets pagination headers.
func (h *ConstructionDetailsHandler) GetAllConstructionDetails(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var propertyID *uuid.UUID
	if propertyIDStr := c.Query("propertyId"); propertyIDStr != "" {
		if id, err := uuid.Parse(propertyIDStr); err == nil {
			propertyID = &id
		}
	}

	constructionDetails, total, err := h.constructionDetailsService.GetAllConstructionDetails(c.Request.Context(), page, size, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get construction details", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, constructionDetails)
}

// GetConstructionDetailsByPropertyID handles GET requests to retrieve construction details by property ID.
func (h *ConstructionDetailsHandler) GetConstructionDetailsByPropertyID(c *gin.Context) {
	propertyIDStr := c.Param("propertyId")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	constructionDetails, err := h.constructionDetailsService.GetConstructionDetailsByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get construction details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Construction details retrieved successfully", constructionDetails))
}
