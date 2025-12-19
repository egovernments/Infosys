package handlers

import (
    "enumeration/internal/models"
    "enumeration/internal/services"
    "enumeration/pkg/response"
    "strconv"
    "time"
     "github.com/lib/pq"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

// AmenityRequest represents the request body for creating/updating amenities
type AmenityRequest struct {
    PropertyID  string     `json:"property_id" binding:"required,uuid"`
    Type        pq.StringArray  `gorm:"type:text[]" json:"type" binding:"required,dive"`
    Description string     `json:"description"`
    ExpiryDate  *time.Time `json:"expiry_date"`
}

type AmenityHandler struct {
    service services.AmenityService
}

func NewAmenityHandler(service services.AmenityService) *AmenityHandler {
    return &AmenityHandler{service}
}

func (h *AmenityHandler) GetAll(c *gin.Context) {
    // Parse query parameters; service will enforce defaults/bounds
    page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
    amenityType := c.Query("type")
    propertyID := c.Query("propertyId")

    var (
        amenities []models.Amenities
        total     int64
        err       error
    )

    // If no filters provided, call plain paginated GetAll; otherwise call filtered paginated method
    if amenityType == "" && propertyID == "" {
        amenities, total, err = h.service.GetAll(c.Request.Context(), page, size)
    } else {
        amenities, total, err = h.service.GetAllWithFilters(c.Request.Context(), page, size, amenityType, propertyID)
    }
    if err != nil {
        response.InternalServerError(c, "Failed to retrieve amenities: "+err.Error())
        return
    }

    // Calculate pagination meta (use the requested 'size' to compute pages; service has already clamped size)
    totalPages := int((total + int64(size) - 1) / int64(size))

    // Return paginated response
    response.Paginated(c, amenities, response.PaginationMeta{
        Page:       page,
        PageSize:   size,
        TotalItems: total,
        TotalPages: totalPages,
    })
}

func (h *AmenityHandler) GetByID(c *gin.Context) {
    id := c.Param("id")

    // Validate UUID format
    if _, err := uuid.Parse(id); err != nil {
        response.BadRequest(c, "Invalid ID format")
        return
    }

    amenity, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.NotFound(c, "Amenity not found")
        return
    }

    response.Success(c, "Amenity retrieved successfully", amenity)
}

func (h *AmenityHandler) Create(c *gin.Context) {
    var req AmenityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request body: "+err.Error())
        return
    }

    // Parse property ID
    propertyUUID, err := uuid.Parse(req.PropertyID)
    if err != nil {
        response.BadRequest(c, "Invalid property ID format")
        return
    }

    // Create amenity model
    amenity := &models.Amenities{
        PropertyID:  propertyUUID,
        Type:        req.Type,
        Description: req.Description,
        ExpiryDate:  req.ExpiryDate,
    }

    if err := h.service.Create(c.Request.Context(), amenity); err != nil {
        response.InternalServerError(c, "Failed to create amenity: "+err.Error())
        return
    }

    response.Created(c, map[string]interface{}{
        "message": "Amenity created successfully",
        "data":    amenity,
    })
}

func (h *AmenityHandler) Update(c *gin.Context) {
    id := c.Param("id")

    // Validate UUID format
    if _, err := uuid.Parse(id); err != nil {
        response.BadRequest(c, "Invalid ID format")
        return
    }

    var req AmenityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request body: "+err.Error())
        return
    }

    // Parse property ID
    propertyUUID, err := uuid.Parse(req.PropertyID)
    if err != nil {
        response.BadRequest(c, "Invalid property ID format")
        return
    }

    // Create amenity model for update
    amenity := &models.Amenities{
        PropertyID:  propertyUUID,
        Type:        req.Type,
        Description: req.Description,
        ExpiryDate:  req.ExpiryDate,
    }

    if err := h.service.Update(c.Request.Context(), id, amenity); err != nil {
        response.InternalServerError(c, "Failed to update amenity: "+err.Error())
        return
    }

    // Get updated amenity to return
    updatedAmenity, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.InternalServerError(c, "Failed to retrieve updated amenity")
        return
    }

    response.Success(c, "Amenity updated successfully", updatedAmenity)
}

func (h *AmenityHandler) Delete(c *gin.Context) {
    id := c.Param("id")

    // Validate UUID format
    if _, err := uuid.Parse(id); err != nil {
        response.BadRequest(c, "Invalid ID format")
        return
    }

    // Check if amenity exists
    _, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.NotFound(c, "Amenity not found")
        return
    }

    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        response.InternalServerError(c, "Failed to delete amenity: "+err.Error())
        return
    }

    response.Success(c, "Amenity deleted successfully", nil)
}

func (h *AmenityHandler) GetByPropertyID(c *gin.Context) {
    propertyId := c.Param("propertyId")

    // Validate UUID format
    if _, err := uuid.Parse(propertyId); err != nil {
        response.BadRequest(c, "Invalid property ID format")
        return
    }

    amenities, err := h.service.GetByPropertyID(c.Request.Context(), propertyId)
    if err != nil {
        response.InternalServerError(c, "Failed to retrieve amenities by property ID: "+err.Error())
        return
    }

    response.Success(c, "Amenities retrieved successfully", amenities)
}