// Package response provides standardized response structures and helper functions
// for sending JSON responses in a RESTful API using the Gin framework.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response structure for API errors.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse represents a standard success response structure for API responses.
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationMeta contains metadata for paginated API responses.
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

// PaginatedResponse represents a paginated API response with data and pagination info.
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// Success sends a 200 OK JSON response with a message and optional data.
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created JSON response with the provided data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// BadRequest sends a 400 Bad Request error response with a message.
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "BadRequest",
		Message: message,
		Code:    http.StatusBadRequest,
	})
}

// NotFound sends a 404 Not Found error response with a message.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "NotFound",
		Message: message,
		Code:    http.StatusNotFound,
	})
}

// Unauthorized sends a 401 Unauthorized error response with a message.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error:   "Unauthorized",
		Message: message,
		Code:    http.StatusUnauthorized,
	})
}

// Forbidden sends a 403 Forbidden error response with a message.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Error:   "Forbidden",
		Message: message,
		Code:    http.StatusForbidden,
	})
}

// InternalServerError sends a 500 Internal Server Error response with a message.
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "InternalError",
		Message: message,
		Code:    http.StatusInternalServerError,
	})
}

// CustomError sends an error response with a custom status code, error type, and message.
func CustomError(c *gin.Context, code int, errorType string, message string) {
	c.JSON(code, ErrorResponse{
		Error:   errorType,
		Message: message,
		Code:    code,
	})
}

// Paginated sends a 200 OK paginated response with data and pagination metadata.
func Paginated(c *gin.Context, data interface{}, pagination PaginationMeta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	})
}

// ApiResponse is a generic API response structure supporting both success and error cases.
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// SuccessResponseBody creates a success ApiResponse with a message and optional data.
func SuccessResponseBody(message string, data interface{}) ApiResponse {
	return ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponseBody creates an error ApiResponse with a message and optional error details.
func ErrorResponseBody(message string, errors ...string) ApiResponse {
	return ApiResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
