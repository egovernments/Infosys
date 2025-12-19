// Package handlers contains HTTP handler implementations for the property tax enumeration system.
package handlers

import (
	"enumeration/internal/models"
	"enumeration/internal/services"
	"enumeration/pkg/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AssessmentDetailsHandler handles HTTP requests for assessment details resources.
type AssessmentDetailsHandler struct {
	assessmentDetailsService services.AssessmentDetailsService // Service layer for assessment details operations
}

// NewAssessmentDetailsHandler creates a new AssessmentDetailsHandler with the provided service.
func NewAssessmentDetailsHandler(assessmentDetailsService services.AssessmentDetailsService) *AssessmentDetailsHandler {
	return &AssessmentDetailsHandler{
		assessmentDetailsService: assessmentDetailsService,
	}
}

// CreateAssessmentDetails handles POST requests to create a new assessment details record.
// Validates the request body and delegates creation to the service layer.
func (h *AssessmentDetailsHandler) CreateAssessmentDetails(c *gin.Context) {
	var assessmentDetails models.AssessmentDetails

	if err := c.ShouldBindJSON(&assessmentDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	if err := h.assessmentDetailsService.CreateAssessmentDetails(c.Request.Context(), &assessmentDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to create assessment details", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponseBody("Assessment details created successfully", assessmentDetails))
}

// GetAssessmentDetailsByID handles GET requests to retrieve assessment details by their ID.
func (h *AssessmentDetailsHandler) GetAssessmentDetailsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid assessment details ID", err.Error()))
		return
	}

	assessmentDetails, err := h.assessmentDetailsService.GetAssessmentDetailsByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "assessment details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Assessment details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get assessment details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Assessment details retrieved successfully", assessmentDetails))
}

// UpdateAssessmentDetails handles PUT requests to update existing assessment details by their ID.
// Validates the request body and delegates update to the service layer.
func (h *AssessmentDetailsHandler) UpdateAssessmentDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid assessment details ID", err.Error()))
		return
	}

	var assessmentDetails models.AssessmentDetails
	if err := c.ShouldBindJSON(&assessmentDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	assessmentDetails.ID = id

	if err := h.assessmentDetailsService.UpdateAssessmentDetails(c.Request.Context(), &assessmentDetails); err != nil {
		if err.Error() == "assessment details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Assessment details not found", err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to update assessment details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Assessment details updated successfully", assessmentDetails))
}

// DeleteAssessmentDetails handles DELETE requests to remove assessment details by their ID.
func (h *AssessmentDetailsHandler) DeleteAssessmentDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid assessment details ID", err.Error()))
		return
	}

	if err := h.assessmentDetailsService.DeleteAssessmentDetails(c.Request.Context(), id); err != nil {
		if err.Error() == "assessment details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Assessment details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to delete assessment details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Assessment details deleted successfully", nil))
}

// GetAllAssessmentDetails handles GET requests to retrieve all assessment details with pagination.
// Optionally filters by property ID and sets pagination headers.
func (h *AssessmentDetailsHandler) GetAllAssessmentDetails(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var propertyID *uuid.UUID
	if propertyIDStr := c.Query("propertyId"); propertyIDStr != "" {
		if id, err := uuid.Parse(propertyIDStr); err == nil {
			propertyID = &id
		}
	}

	assessmentDetails, total, err := h.assessmentDetailsService.GetAllAssessmentDetails(c.Request.Context(), page, size, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get assessment details", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, assessmentDetails)
}

// GetAssessmentDetailsByPropertyID handles GET requests to retrieve assessment details by property ID.
func (h *AssessmentDetailsHandler) GetAssessmentDetailsByPropertyID(c *gin.Context) {
	propertyIDStr := c.Param("propertyId")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	assessmentDetails, err := h.assessmentDetailsService.GetAssessmentDetailsByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Assessment details not found for this property", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get assessment details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Assessment details retrieved successfully", assessmentDetails))
}
