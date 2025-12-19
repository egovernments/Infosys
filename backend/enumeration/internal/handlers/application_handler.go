// Package handlers contains HTTP handlers for API endpoints
package handlers

import (
	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/services"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApplicationHandler handles HTTP requests for applications
type ApplicationHandler struct {
	service services.ApplicationService
}

// NewApplicationHandler creates a new ApplicationHandler instance
func NewApplicationHandler(service services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

// Create handles POST /applications to create a new application
func (h *ApplicationHandler) Create(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}
	citizenID := c.GetHeader(constants.HeaderUserID)
	tenantID := c.GetHeader(constants.HeaderTenantID)
	application, err := h.service.Create(c.Request.Context(), tenantID, citizenID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": constants.ErrApplicationCreationFailed,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Application created successfully",
		"data":    application,
	})
}

// GetByID handles GET /applications/:id to fetch an application by ID
func (h *ApplicationHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	application, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": constants.ErrApplicationNotFound,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application retrieved successfully",
		"data":    application,
	})
}

// GetByApplicationNo handles GET /applications/by-number/:applicationNo to fetch by application number
func (h *ApplicationHandler) GetByApplicationNo(c *gin.Context) {
	applicationNo := c.Param("applicationNo")
	if applicationNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Application number is required",
			"errors":  []string{},
		})
		return
	}

	application, err := h.service.GetByApplicationNo(c.Request.Context(), applicationNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": constants.ErrApplicationNotFound,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application retrieved successfully",
		"data":    application,
	})
}

// Update handles PUT /applications/:id to update an application
func (h *ApplicationHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	var req dto.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	application, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": constants.ErrApplicationUpdateFailed,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application updated successfully",
		"data":    application,
	})
}

// UpdateStatus handles PATCH /applications/:id/status to update application status
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to update status",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status updated successfully",
	})
}

// Action handles PATCH /applications/:id/:action for assign, verify, audit-verify, approve
func (h *ApplicationHandler) Action(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application ID",
			"errors":  []string{err.Error()},
		})
		return
	}
	tenantID := c.GetHeader(constants.HeaderTenantID)
	userId := c.GetHeader(constants.HeaderUserID)
	var req dto.ActionRequest
	reqInterface, exists := c.Get("actionRequest")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"errors":  []string{},
		})
		return
	}
	req = reqInterface.(dto.ActionRequest)
	fmt.Printf("req: %v\n", req)

	switch req.Action {
	case "assign":
		if err := h.service.AssignAgent(c.Request.Context(), id, &req, tenantID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to assign agent",
				"errors":  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Agent assigned successfully",
		})

	case "verify":
		err := h.service.VerifyApplicationByAgent(c.Request.Context(), tenantID, userId, id.String(), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to verify application",
				"errors":  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Application verified successfully by agent",
		})

	case "audit-verify":
		err := h.service.VerifyApplication(c.Request.Context(), tenantID, userId, id.String(), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to audit verify application",
				"errors":  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Application audit verified successfully by service manager",
		})

	case "approve":
		err := h.service.ApproveApplication(c.Request.Context(), tenantID, userId, id.String(), &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to approve the  application",
				"errors":  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Application approved successfully by commissioner",
		})

	case "re-assign":
		// service manager reassigns an agent
		if err := h.service.AssignAgent(c.Request.Context(), id, &req, tenantID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to reassign agent",
				"errors":  []string{err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Agent reassigned successfully",
		})
	}
}

// Delete handles DELETE /applications/:id to remove an application
func (h *ApplicationHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application ID",
			"errors":  []string{err.Error()},
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": constants.ErrApplicationDeletionFailed,
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application deleted successfully",
	})
}

// List handles GET /applications to list all applications with pagination
func (h *ApplicationHandler) List(c *gin.Context) {

	userID := c.GetHeader(constants.HeaderUserID)
	role := c.GetHeader(constants.HeaderUserRole)
	verify := c.GetHeader(constants.HeaderStatus)
	tenantID := c.GetHeader(constants.HeaderTenantID)

	fmt.Printf("\"in handler\": %v\n", "in handler")
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	applications, total, err := h.service.List(c.Request.Context(), userID, tenantID, role, verify, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve applications",
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
	c.Header(constants.HeaderTotalCount, strconv.FormatInt(total, 10))
	c.Header(constants.HeaderCurrentPage, strconv.Itoa(page))
	c.Header(constants.HeaderPerPage, strconv.Itoa(size))
	c.Header(constants.HeaderTotalPages, strconv.Itoa(totalPages))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Applications retrieved successfully",
		"data":    applications,
		"pagination": gin.H{
			"page":       page,
			"size":       size,
			"totalItems": total,
			"totalPages": totalPages,
		},
	})
}

// Search handles GET /applications/search to search applications with filters
func (h *ApplicationHandler) Search(c *gin.Context) {
	var criteria dto.ApplicationSearchCriteria
	if err := c.ShouldBindQuery(&criteria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid search criteria",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	applications, total, err := h.service.Search(c.Request.Context(), &criteria, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to search applications",
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
	c.Header(constants.HeaderTotalCount, strconv.FormatInt(total, 10))
	c.Header(constants.HeaderCurrentPage, strconv.Itoa(page))
	c.Header(constants.HeaderPerPage, strconv.Itoa(size))
	c.Header(constants.HeaderTotalPages, strconv.Itoa(totalPages))
    if criteria.IsCountOnly {
		c.JSON(http.StatusOK, gin.H{
			"totalItems": total,
			
		},
	)
	}else{
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Applications retrieved successfully",
		"data":    applications,
		"pagination": gin.H{
			"page":       page,
			"size":       size,
			"totalItems": total,
			"totalPages": totalPages,
		},
	})
}
}
