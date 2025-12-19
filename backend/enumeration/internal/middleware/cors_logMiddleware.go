package middleware

import (
	"bytes"
	"enumeration/pkg/logger"
	"io"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware sets the necessary headers to handle CORS (Cross-Origin Resource Sharing) requests.
// It allows any origin, credentials, a set of headers, and methods, and handles OPTIONS requests by returning 204.
func CorsMiddleware() gin.HandlerFunc {
	return (func(c *gin.Context) {
		// Allow all origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow credentials such as cookies, authorization headers, or TLS client certificates
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,X-User-Role,X-User-ID,X-Status,X-Tenant-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Handle preflight OPTIONS request by aborting with status code 204 (No Content)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

// LoggerMiddleware logs the raw request body before passing control to the next handler.
// It restores the request body for subsequent handlers, as reading the body consumes it.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the raw request body into bodyBytes
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			// Log an error if reading fails
			logger.Error("Failed to read request body:", err)
		} else {
			// Log the request body contents
			logger.Info("Request Body:", string(bodyBytes))
			// Restore the request body for the next handlers, as reading it consumes the body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		// Continue with the next handler in the chain
		c.Next()
	}
}
