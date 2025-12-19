package handlers

import (
	"bytes"
	"encoding/json"
	"enumeration/internal/models"
	"enumeration/internal/services"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CoordinatesHandler handles HTTP requests for coordinates resources.
type CoordinatesHandler struct {
	service services.CoordinatesService // Service layer for coordinates operations
}

// NewCoordinatesHandler creates a new CoordinatesHandler with the provided service.
func NewCoordinatesHandler(service services.CoordinatesService) *CoordinatesHandler {
	return &CoordinatesHandler{service: service}
}

// GetAll handles GET /coordinates and returns a paginated list of coordinates.
// Optionally filters by GIS Data ID and sets pagination headers.
func (h *CoordinatesHandler) GetAll(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// Parse optional filter
	var gisDataID *uuid.UUID
	if gisDataIDStr := c.Query("gisDataId"); gisDataIDStr != "" {
		id, err := uuid.Parse(gisDataIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid gisDataId format",
				"errors":  []string{err.Error()},
			})
			return
		}
		gisDataID = &id
	}

	coordinates, total, err := h.service.GetAll(c.Request.Context(), page, size, gisDataID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve coordinates",
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
		"message": "Coordinates retrieved successfully",
		"data":    coordinates,
		"pagination": gin.H{
			"page":       page,
			"size":       size,
			"totalItems": total,
			"totalPages": totalPages,
		},
	})
}

// Create handles POST /coordinates and creates a new coordinates record.
// Validates the request body and delegates creation to the service layer.
func (h *CoordinatesHandler) Create(c *gin.Context) {
	var coordinates models.Coordinates
	if err := c.ShouldBindJSON(&coordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &coordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to create coordinates",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Coordinates created successfully",
		"data":    coordinates,
	})
}

// GetByID handles GET /coordinates/:id and retrieves a coordinates record by its ID.
func (h *CoordinatesHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid coordinates ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	coordinates, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Coordinates not found",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Coordinates retrieved successfully",
		"data":    coordinates,
	})
}

// Update handles PUT /coordinates/:id and updates an existing coordinates record by its ID.
// Validates the request body and delegates update to the service layer.
func (h *CoordinatesHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid coordinates ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	var coordinates models.Coordinates
	if err := c.ShouldBindJSON(&coordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	coordinates.ID = id
	if err := h.service.Update(c.Request.Context(), &coordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to update coordinates",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Coordinates updated successfully",
		"data":    coordinates,
	})
}

// Delete handles DELETE /coordinates/:id and removes a coordinates record by its ID.
func (h *CoordinatesHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid coordinates ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Failed to delete coordinates",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Coordinates deleted successfully",
	})
}

// CreateBatch handles POST /coordinates/batch and accepts either a single object or an array.
// Tries to unmarshal the request body as a slice first, then as a single object if that fails.
func (h *CoordinatesHandler) CreateBatch(c *gin.Context) {
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	// Try unmarshalling into slice first
	var coordsSlice []*models.Coordinates
	if err := json.Unmarshal(raw, &coordsSlice); err == nil {
		// slice succeeded
		if len(coordsSlice) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no coordinates provided"})
			return
		}
		if err := h.service.CreateBatch(c.Request.Context(), coordsSlice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create coordinates", "error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "coordinates created successfully", "data": coordsSlice})
		return
	}

	// If not a slice, try single object
	var single models.Coordinates
	if err := json.Unmarshal(raw, &single); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request payload", "error": err.Error()})
		return
	}
	if err := h.service.Create(c.Request.Context(), &single); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create coordinates", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "coordinate created successfully", "data": single})
}

// ReplaceByGISDataID handles PUT /coordinates/gis/:gisDataId and replaces all coordinates for a given GIS Data ID.
// Accepts either an array or a single object in the request body. An empty body deletes all coordinates for the GIS Data ID.
func (h *CoordinatesHandler) ReplaceByGISDataID(c *gin.Context) {
	gisIDStr := c.Param("gisDataId")
	gisID, err := uuid.Parse(gisIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid gisDataId",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Read raw body to support either array or single object
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body", "errors": []string{err.Error()}})
		return
	}
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	var coordsSlice []*models.Coordinates

	if len(trimmed) == 0 {
		// Empty body interpreted as delete all coordinates for this GISData
		coordsSlice = []*models.Coordinates{}
	} else if trimmed[0] == '[' {
		if err := json.Unmarshal(raw, &coordsSlice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body", "errors": []string{err.Error()}})
			return
		}
	} else {
		var single models.Coordinates
		if err := json.Unmarshal(raw, &single); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request body", "errors": []string{err.Error()}})
			return
		}
		coordsSlice = append(coordsSlice, &single)
	}

	// call service (service sets GISDataID on each coord and validates)
	if err := h.service.ReplaceByGISDataID(c.Request.Context(), gisID, coordsSlice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to replace coordinates", "errors": []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Coordinates replaced successfully",
		"data":    coordsSlice,
	})
}
