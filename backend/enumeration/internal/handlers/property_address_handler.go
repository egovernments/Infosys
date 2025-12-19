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

// PropertyAddressHandler handles HTTP requests for property address resources.
type PropertyAddressHandler struct {
	propertyAddressService services.PropertyAddressService // Service layer for property address operations
}

// NewPropertyAddressHandler creates a new PropertyAddressHandler with the provided service.
func NewPropertyAddressHandler(propertyAddressService services.PropertyAddressService) *PropertyAddressHandler {
	return &PropertyAddressHandler{
		propertyAddressService: propertyAddressService,
	}
}

// CreatePropertyAddress handles POST requests to create a new property address record.
// Validates the request body and delegates creation to the service layer.
func (h *PropertyAddressHandler) CreatePropertyAddress(c *gin.Context) {
	var address models.PropertyAddress

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	if err := h.propertyAddressService.CreatePropertyAddress(c.Request.Context(), &address); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to create property address", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponseBody("Property address created successfully", address))
}

// GetPropertyAddressByID handles GET requests to retrieve a property address by its ID.
func (h *PropertyAddressHandler) GetPropertyAddressByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property address ID", err.Error()))
		return
	}

	address, err := h.propertyAddressService.GetPropertyAddressByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "property address not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Property address not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get property address", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property address retrieved successfully", address))
}

// UpdatePropertyAddress handles PUT requests to update an existing property address by its ID.
// Validates the request body and delegates update to the service layer.
func (h *PropertyAddressHandler) UpdatePropertyAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property address ID", err.Error()))
		return
	}

	var address models.PropertyAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	address.ID = id

	if err := h.propertyAddressService.UpdatePropertyAddress(c.Request.Context(), &address); err != nil {
		if err.Error() == "property address not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Property address not found", err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to update property address", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property address updated successfully", address))
}

// DeletePropertyAddress handles DELETE requests to remove a property address by its ID.
func (h *PropertyAddressHandler) DeletePropertyAddress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property address ID", err.Error()))
		return
	}

	if err := h.propertyAddressService.DeletePropertyAddress(c.Request.Context(), id); err != nil {
		if err.Error() == "property address not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Property address not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to delete property address", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property address deleted successfully", nil))
}

// GetAllPropertyAddresses handles GET requests to retrieve all property addresses with pagination and optional filtering by property ID.
// Sets pagination headers in the response.
func (h *PropertyAddressHandler) GetAllPropertyAddresses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var propertyID *uuid.UUID
	if propertyIDStr := c.Query("propertyId"); propertyIDStr != "" {
		if id, err := uuid.Parse(propertyIDStr); err == nil {
			propertyID = &id
		}
	}

	addresses, total, err := h.propertyAddressService.GetAllPropertyAddresses(c.Request.Context(), page, size, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get property addresses", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, addresses)
}

// GetPropertyAddressByPropertyID handles GET requests to retrieve a property address by property ID.
func (h *PropertyAddressHandler) GetPropertyAddressByPropertyID(c *gin.Context) {
	propertyIDStr := c.Param("propertyId")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid property ID", err.Error()))
		return
	}

	address, err := h.propertyAddressService.GetPropertyAddressByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		if err.Error() == "property address not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Property address not found for this property", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get property address", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Property address retrieved successfully", address))
}

// SearchPropertyAddresses handles GET requests to search property addresses with advanced filters and pagination.
// Supports filtering by property ID, locality, zone, ward, block, street, election ward, secretariat ward, and pin code.
// Sets pagination headers in the response.
func (h *PropertyAddressHandler) SearchPropertyAddresses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	sortBy := c.DefaultQuery("sortBy", "createdAt")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	params := services.SearchPropertyAddressParams{
		Page:      page,
		Size:      size,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	// Optional filters
	if propertyID := c.Query("propertyId"); propertyID != "" {
		if id, err := uuid.Parse(propertyID); err == nil {
			params.PropertyID = &id
		}
	}

	if locality := c.Query("locality"); locality != "" {
		params.Locality = &locality
	}

	if zoneNo := c.Query("zoneNo"); zoneNo != "" {
		params.ZoneNo = &zoneNo
	}

	if wardNo := c.Query("wardNo"); wardNo != "" {
		params.WardNo = &wardNo
	}

	if blockNo := c.Query("blockNo"); blockNo != "" {
		params.BlockNo = &blockNo
	}

	if street := c.Query("street"); street != "" {
		params.Street = &street
	}

	if electionWard := c.Query("electionWard"); electionWard != "" {
		params.ElectionWard = &electionWard
	}

	if secretariatWard := c.Query("secretariatWard"); secretariatWard != "" {
		params.SecretariatWard = &secretariatWard
	}

	if pinCodeStr := c.Query("pinCode"); pinCodeStr != "" {
		if pinCode, err := strconv.ParseUint(pinCodeStr, 10, 64); err == nil {
			params.PinCode = &pinCode
		}
	}

	addresses, total, err := h.propertyAddressService.SearchPropertyAddresses(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to search property addresses", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, addresses)
}
