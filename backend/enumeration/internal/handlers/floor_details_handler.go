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

// FloorDetailsHandler handles HTTP requests for floor details resources.
type FloorDetailsHandler struct {
	floorDetailsService services.FloorDetailsService // Service layer for floor details operations
}

// NewFloorDetailsHandler creates a new FloorDetailsHandler with the provided service.
func NewFloorDetailsHandler(floorDetailsService services.FloorDetailsService) *FloorDetailsHandler {
	return &FloorDetailsHandler{
		floorDetailsService: floorDetailsService,
	}
}

// CreateFloorDetails handles POST requests to create a new floor details record.
// Validates the request body and delegates creation to the service layer.
func (h *FloorDetailsHandler) CreateFloorDetails(c *gin.Context) {
	var floorDetails models.FloorDetails

	if err := c.ShouldBindJSON(&floorDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	if err := h.floorDetailsService.CreateFloorDetails(c.Request.Context(), &floorDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to create floor details", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponseBody("Floor details created successfully", floorDetails))
}

// GetFloorDetailsByID handles GET requests to retrieve floor details by their ID.
func (h *FloorDetailsHandler) GetFloorDetailsByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid floor details ID", err.Error()))
		return
	}

	floorDetails, err := h.floorDetailsService.GetFloorDetailsByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "floor details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Floor details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get floor details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Floor details retrieved successfully", floorDetails))
}

// UpdateFloorDetails handles PUT requests to update existing floor details by their ID.
// Validates the request body and delegates update to the service layer.
func (h *FloorDetailsHandler) UpdateFloorDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid floor details ID", err.Error()))
		return
	}

	var floorDetails models.FloorDetails
	if err := c.ShouldBindJSON(&floorDetails); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid request body", err.Error()))
		return
	}

	floorDetails.ID = id

	if err := h.floorDetailsService.UpdateFloorDetails(c.Request.Context(), &floorDetails); err != nil {
		if err.Error() == "floor details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Floor details not found", err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Failed to update floor details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Floor details updated successfully", floorDetails))
}

// DeleteFloorDetails handles DELETE requests to remove floor details by their ID.
func (h *FloorDetailsHandler) DeleteFloorDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponseBody("Invalid floor details ID", err.Error()))
		return
	}

	if err := h.floorDetailsService.DeleteFloorDetails(c.Request.Context(), id); err != nil {
		if err.Error() == "floor details not found" {
			c.JSON(http.StatusNotFound, response.ErrorResponseBody("Floor details not found", err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to delete floor details", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponseBody("Floor details deleted successfully", nil))
}

// GetAllFloorDetails handles GET requests to retrieve all floor details with pagination.
// Optionally filters by construction details ID and sets pagination headers.
func (h *FloorDetailsHandler) GetAllFloorDetails(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var constructionDetailsID *uuid.UUID
	if constructionDetailsIDStr := c.Query("constructionDetailsId"); constructionDetailsIDStr != "" {
		if id, err := uuid.Parse(constructionDetailsIDStr); err == nil {
			constructionDetailsID = &id
		}
	}

	floorDetails, total, err := h.floorDetailsService.GetAllFloorDetails(c.Request.Context(), page, size, constructionDetailsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponseBody("Failed to get floor details", err.Error()))
		return
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))

	c.JSON(http.StatusOK, floorDetails)
}
