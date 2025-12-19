package handlers

import (
	"net/http"
	"strconv"

	"enumeration/internal/constants"
	"enumeration/internal/models"
	"enumeration/internal/services"
	"enumeration/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PropertyHandler handles HTTP requests for property resources.
type PropertyHandler struct {
	propertyService services.PropertyService // Service layer for property operations
}

// NewPropertyHandler creates a new PropertyHandler with the provided service.
func NewPropertyHandler(propertyService services.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		propertyService: propertyService,
	}
}

// CreateProperty handles POST requests to create a new property record.
// Validates the request body and delegates creation to the service layer.
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	var property models.Property

	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	if err := h.propertyService.CreateProperty(c.Request.Context(), &property); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to create property", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponseBody("Property created successfully", property))
}

// GetPropertyByID handles GET requests to retrieve a property by its ID.
func (h *PropertyHandler) GetPropertyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	property, err := h.propertyService.GetPropertyByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == constants.ErrPropertyNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody(constants.ErrPropertyNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get property", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property retrieved successfully", property))
}

// UpdateProperty handles PUT requests to update an existing property by its ID.
// Validates the request body and delegates update to the service layer.
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	property.ID = id

	if err := h.propertyService.UpdateProperty(c.Request.Context(), &property); err != nil {
		
			c.JSON(http.StatusNotFound, response.ErrorResponseBody(err.Error()))
			return
		
	}
	updatedProperty, err := h.propertyService.GetPropertyByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to retrieve updated property", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property updated successfully", updatedProperty))
}

// DeleteProperty handles DELETE requests to remove a property by its ID.
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	if err := h.propertyService.DeleteProperty(c.Request.Context(), id); err != nil {
		if err.Error() == constants.ErrPropertyNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody(constants.ErrPropertyNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to delete property", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property deleted successfully", nil))
}

// GetAllProperties handles GET requests to retrieve all properties with pagination and optional filtering by property type.
// Sets pagination headers in the response.
func (h *PropertyHandler) GetAllProperties(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var propertyType *string
	if propertyTypeStr := c.Query("propertyType"); propertyTypeStr != "" {
		propertyType = &propertyTypeStr
	}

	properties, total, err := h.propertyService.GetAllProperties(c.Request.Context(), page, size, propertyType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get properties", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, properties)
}

// GetPropertyByPropertyNo handles GET requests to retrieve a property by its property number.
func (h *PropertyHandler) GetPropertyByPropertyNo(c *gin.Context) {
	propertyNo := c.Param("propertyNo")
	if propertyNo == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Property number is required", ""))
		return
	}

	property, err := h.propertyService.GetPropertyByPropertyNo(c.Request.Context(), propertyNo)
	if err != nil {
		if err.Error() == constants.ErrPropertyNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody(constants.ErrPropertyNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get property", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property retrieved successfully", property))
}

// SearchProperties handles GET requests to search properties with advanced filters and pagination.
// Supports filtering by property type, ownership type, complex name, locality, ward, zone, and street.
// Sets pagination headers in the response.
func (h *PropertyHandler) SearchProperties(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	sortBy := c.DefaultQuery("sortBy", "createdAt")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	params := services.SearchPropertyParams{
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// Optional filters
	if propertyType := c.Query("propertyType"); propertyType != "" {
		params.PropertyType = &propertyType
	}

	if ownershipType := c.Query("ownershipType"); ownershipType != "" {
		params.OwnershipType = &ownershipType
	}

	if complexName := c.Query("complexName"); complexName != "" {
		params.ComplexName = &complexName
	}

	if locality := c.Query("locality"); locality != "" {
		params.Locality = &locality
	}

	if wardNo := c.Query("wardNo"); wardNo != "" {
		params.WardNo = &wardNo
	}

	if zoneNo := c.Query("zoneNo"); zoneNo != "" {
		params.ZoneNo = &zoneNo
	}

	if street := c.Query("street"); street != "" {
		params.Street = &street
	}

	properties, total, err := h.propertyService.SearchProperties(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to search properties", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, properties)
}
