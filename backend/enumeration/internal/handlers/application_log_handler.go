package handlers

import (
	"enumeration/internal/dto"
	"enumeration/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationLogHandler struct {
	service services.ApplicationLogService
}

func NewApplicationLogHandler(service services.ApplicationLogService) *ApplicationLogHandler {
	return &ApplicationLogHandler{
		service: service,
	}
}

// GetApplicationLogs retrieves all application logs with pagination and filtering
func (h *ApplicationLogHandler) GetApplicationLogs(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	applicationID := ctx.Query("applicationId")
	action := ctx.Query("action")

	logs, total, err := h.service.List(ctx, applicationID, action, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.Header("X-Current-Page", strconv.Itoa(page))
	ctx.Header("X-Per-Page", strconv.Itoa(size))

	ctx.JSON(http.StatusOK, logs)
}

// CreateApplicationLog creates a new application log entry
func (h *ApplicationLogHandler) CreateApplicationLog(ctx *gin.Context) {
	var req dto.CreateApplicationLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.service.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create application log",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Application log created successfully",
		"data":    log,
	})
}

// GetApplicationLogByID retrieves a specific application log by ID
func (h *ApplicationLogHandler) GetApplicationLogByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	log, err := h.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Application log not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    log,
	})
}

// UpdateApplicationLog updates an existing application log
func (h *ApplicationLogHandler) UpdateApplicationLog(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	var req dto.UpdateApplicationLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.service.Update(ctx, id, &req)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Application log not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update application log",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application log updated successfully",
		"data":    log,
	})
}

// DeleteApplicationLog deletes an application log by ID
func (h *ApplicationLogHandler) DeleteApplicationLog(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log ID"})
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Application log not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete application log",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application log deleted successfully",
	})
}

// GetApplicationLogsByApplicationID retrieves all logs for a specific application
func (h *ApplicationLogHandler) GetApplicationLogsByApplicationID(ctx *gin.Context) {
	applicationID, err := uuid.Parse(ctx.Param("applicationId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	action := ctx.Query("action")

	logs, total, err := h.service.GetByApplicationID(ctx, applicationID, action, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.Header("X-Current-Page", strconv.Itoa(page))
	ctx.Header("X-Per-Page", strconv.Itoa(size))

	ctx.JSON(http.StatusOK, logs)
}
