package handlers

import (
	"enumeration/internal/models"
	"enumeration/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GISHandler handles HTTP requests for GIS data resources.
type GISHandler struct {
	service services.GISService // Service layer for GIS data operations
}

// NewGISHandler creates a new GISHandler with the provided service.
func NewGISHandler(service services.GISService) *GISHandler {
	return &GISHandler{service: service}
}

// GetAll handles GET /gis-data and returns a paginated list of all GIS data.
// Sets pagination headers in the response.
func (h *GISHandler) GetAll(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	gisDataList, total, err := h.service.GetAllGISData(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve GIS data",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Calculate pagination
	totalPages := int(total) / size
	if int(total)%size != 0 {
		totalPages++
	}

	// Set pagination headers
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(size))
	c.Header("X-Total-Pages", strconv.Itoa(totalPages))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "GIS data retrieved successfully",
		"data":    gisDataList,
		"pagination": gin.H{
			"page":       page,
			"size":       size,
			"totalItems": total,
			"totalPages": totalPages,
		},
	})
}

// Create handles POST /gis-data and creates a new GIS data record.
// Validates the request body and delegates creation to the service layer.
func (h *GISHandler) Create(c *gin.Context) {
	var gisData models.GISData
	if err := c.ShouldBindJSON(&gisData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.service.CreateGISData(c.Request.Context(), &gisData); err != nil {
		// Check if it's a validation error or conflict
		statusCode := http.StatusBadRequest
		if err.Error() == "GIS data already exists for this property" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "Failed to create GIS data",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "GIS data created successfully",
		"data":    gisData,
	})
}

// GetByID handles GET /gis-data/{id} and retrieves a GIS data record by its ID.
func (h *GISHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid GIS data ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	gisData, err := h.service.GetGISDataByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "GIS data not found",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "GIS data retrieved successfully",
		"data":    gisData,
	})
}

// Update handles PUT /gis-data/{id} and updates an existing GIS data record by its ID.
// Validates the request body and delegates update to the service layer.
func (h *GISHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid GIS data ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	var gisData models.GISData
	if err := c.ShouldBindJSON(&gisData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Set the ID from the URL parameter
	gisData.ID = id

	if err := h.service.UpdateGISData(c.Request.Context(), &gisData); err != nil {
		// Check if it's a not found error
		statusCode := http.StatusBadRequest
		if err.Error() == "GIS data not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "Failed to update GIS data",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "GIS data updated successfully",
		"data":    gisData,
	})
}

// Delete handles DELETE /gis-data/{id} and removes a GIS data record by its ID.
func (h *GISHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid GIS data ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.service.DeleteGISData(c.Request.Context(), id); err != nil {
		// Check if it's a not found error
		statusCode := http.StatusInternalServerError
		if err.Error() == "GIS data not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "Failed to delete GIS data",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "GIS data deleted successfully",
	})
}

// GetByPropertyID handles GET /gis-data/property/{propertyId} and retrieves GIS data by property ID.
func (h *GISHandler) GetByPropertyID(c *gin.Context) {
	propertyID, err := uuid.Parse(c.Param("propertyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid property ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	gisData, err := h.service.GetGISDataByPropertyID(c.Request.Context(), propertyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "GIS data not found for property",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "GIS data retrieved successfully",
		"data":    gisData,
	})
}
