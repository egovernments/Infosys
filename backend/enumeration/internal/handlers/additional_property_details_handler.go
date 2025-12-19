// Package handlers contains HTTP handlers for API endpoints
package handlers

import (
	"encoding/json"
	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/services"
	"enumeration/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AdditionalPropertyDetailsHandler struct {
	service services.AdditionalPropertyDetailsService // Business logic service
}

// NewAdditionalPropertyDetailsHandler creates a new handler instance
func NewAdditionalPropertyDetailsHandler(service services.AdditionalPropertyDetailsService) *AdditionalPropertyDetailsHandler {
	return &AdditionalPropertyDetailsHandler{
		service: service,
	}
}

// CreateAdditionalPropertyDetailsRequest is the request body for creating additional property details
type CreateAdditionalPropertyDetailsRequest struct {
	FieldName  string      `json:"fieldName" binding:"required"`
	FieldValue interface{} `json:"fieldValue" binding:"required"`
	PropertyID uuid.UUID   `json:"propertyId" binding:"required"`
}

// UpdateAdditionalPropertyDetailsRequest is the request body for updating additional property details
type UpdateAdditionalPropertyDetailsRequest struct {
	FieldName  string      `json:"fieldName" binding:"required"`
	FieldValue interface{} `json:"fieldValue" binding:"required"`
	PropertyID uuid.UUID   `json:"propertyId" binding:"required"`
}

// CreateAdditionalPropertyDetails handles POST to create a new additional property details record
func (h *AdditionalPropertyDetailsHandler) CreateAdditionalPropertyDetails(c *gin.Context) {
	var req CreateAdditionalPropertyDetailsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	// Convert interface{} to JSON bytes
	fieldValueBytes, err := json.Marshal(req.FieldValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid field value format", err.Error()))
		return
	}

	details := models.AdditionalPropertyDetails{
		FieldName:  req.FieldName,
		FieldValue: datatypes.JSON(fieldValueBytes),
		PropertyID: req.PropertyID,
	}

	if err := h.service.CreateAdditionalPropertyDetails(c.Request.Context(), &details); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to create additional property details", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponseBody("Additional property details created successfully", details))
}

// GetAdditionalPropertyDetailsByID handles GET by ID for additional property details
func (h *AdditionalPropertyDetailsHandler) GetAdditionalPropertyDetailsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid additional property details ID", err.Error()))
		return
	}

	details, err := h.service.GetAdditionalPropertyDetailsByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "additional property details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Additional property details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get additional property details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Additional property details retrieved successfully", details))
}

// UpdateAdditionalPropertyDetails handles PUT to update existing additional property details
func (h *AdditionalPropertyDetailsHandler) UpdateAdditionalPropertyDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid additional property details ID", err.Error()))
		return
	}

	var req UpdateAdditionalPropertyDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	// Convert interface{} to JSON bytes
	fieldValueBytes, err := json.Marshal(req.FieldValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid field value format", err.Error()))
		return
	}

	details := models.AdditionalPropertyDetails{
		ID:         id,
		FieldName:  req.FieldName,
		FieldValue: datatypes.JSON(fieldValueBytes),
		PropertyID: req.PropertyID,
	}

	if err := h.service.UpdateAdditionalPropertyDetails(c.Request.Context(), &details); err != nil {
		if err.Error() == "additional property details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Additional property details not found", err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to update additional property details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Additional property details updated successfully", details))
}

// DeleteAdditionalPropertyDetails handles DELETE by ID for additional property details
func (h *AdditionalPropertyDetailsHandler) DeleteAdditionalPropertyDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid additional property details ID", err.Error()))
		return
	}

	if err := h.service.DeleteAdditionalPropertyDetails(c.Request.Context(), id); err != nil {
		if err.Error() == "additional property details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Additional property details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to delete additional property details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Additional property details deleted successfully", nil))
}

// GetAllAdditionalPropertyDetails handles GET for all additional property details with pagination and filtering
func (h *AdditionalPropertyDetailsHandler) GetAllAdditionalPropertyDetails(c *gin.Context) {
	page, _ := strconv.Atoi	(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var propertyID *uuid.UUID
	if propertyIDStr := c.Query("propertyId"); propertyIDStr != "" {
		if id, err := uuid.Parse(propertyIDStr); err == nil {
			propertyID = &id
		}
	}

	var fieldName *string
	if fieldNameStr := c.Query("fieldName"); fieldNameStr != "" {
		fieldName = &fieldNameStr
	}

	details, total, err := h.service.GetAllAdditionalPropertyDetails(c.Request.Context(), page, size, propertyID, fieldName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get additional property details", err.Error()))
		return
	}

	// Set pagination headers
	c.Header(constants.HeaderTotalCount, strconv.FormatInt(total, 10))
	c.Header(constants.HeaderCurrentPage, strconv.Itoa(page))
	c.Header(constants.HeaderPerPage, strconv.Itoa(size))
	c.Header(constants.HeaderTotalPages, strconv.Itoa((int((total + int64(size) - 1) / int64(size)))))

	c.JSON(http.StatusOK, details)
}

// GetAdditionalPropertyDetailsByPropertyID handles GET by property ID for additional property details
func (h *AdditionalPropertyDetailsHandler) GetAdditionalPropertyDetailsByPropertyID(c *gin.Context) {
	propertyIDStr := c.Param("propertyId")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	details, err := h.service.GetAdditionalPropertyDetailsByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get additional property details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Additional property details retrieved successfully", details))
}

// GetAdditionalPropertyDetailsByFieldName handles GET by field name for additional property details
func (h *AdditionalPropertyDetailsHandler) GetAdditionalPropertyDetailsByFieldName(c *gin.Context) {
	fieldName := c.Param("fieldName")
	if fieldName == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Field name is required", ""))
		return
	}

	details, err := h.service.GetAdditionalPropertyDetailsByFieldName(c.Request.Context(), fieldName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get additional property details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Additional property details retrieved successfully", details))
}
