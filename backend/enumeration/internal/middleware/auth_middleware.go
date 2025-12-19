// Package middleware provides authentication and authorization middleware for the API
package middleware

import (
	"context"
	"encoding/json"
	"enumeration/internal/config"
	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/security"
	"enumeration/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// MDMSResponse represents the response structure from MDMS API
type MDMSResponse struct {
	MDMS []MDMSData `json:"mdms"`
}

// MDMSData holds individual MDMS API role data
type MDMSData struct {
	ID               string      `json:"id"`
	TenantID         string      `json:"tenantId"`
	SchemaCode       string      `json:"schemaCode"`
	UniqueIdentifier string      `json:"uniqueIdentifier"`
	Data             APIRoleData `json:"data"`
	IsActive         bool        `json:"isActive"`
}

// APIRoleData holds allowed roles for an API endpoint
type APIRoleData struct {
	Method       string   `json:"method"`
	Endpoint     string   `json:"endpoint"`
	AllowedRoles []string `json:"allowedRoles"`
}

// RoleActionData holds actions allowed for a role
type RoleActionData struct {
	Role    string   `json:"role"`
	Actions []string `json:"actions"`
}

// rolePermissions holds role-action permissions loaded from MDMS
var rolePermissions map[string][]string

// mdmsRolePermissions holds endpoint-role permissions loaded from MDMS
var mdmsRolePermissions map[string][]string

// Init loads role and endpoint permissions from MDMS at startup
func Init() {
	// Load role-action permissions from MDMS
	roleActionPermissions, err := loadRoleActionsFromMDMS()
	if err != nil {
		logger.Error("\nFailed to load role-action permissions from MDMS:\n", err)
	} else {
		rolePermissions = roleActionPermissions
		logger.Info("\nSuccessfully loaded role-action permissions from MDMS :\t", rolePermissions, "\n")
	}

	// Load MDMS API roles once at startup
	mdmsRoles, err := loadAllMDMSRoles()
	if err != nil {
		logger.Error("\nFailed to load MDMS API roles:\n", err)
		// Initialize empty map if MDMS fails
		mdmsRolePermissions = make(map[string][]string)
	} else {
		mdmsRolePermissions = mdmsRoles
		logger.Info("\nSuccessfully loaded MDMS API roles :\t", mdmsRolePermissions, "\n")
	}
}

// RoleMiddleware checks if the user has the necessary role to perform the action specified in the request.
func ActionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Keycloak claims from context (set by AuthMiddleware)
		claimsVal, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "MissingClaims"})
			return
		}
		claims, ok := claimsVal.(*security.KeycloakClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "InvalidClaimsType"})
			return
		}

		// Extract roles from Keycloak claims (RealmAccess["roles"])
		var roles []string
		if claims.RealmAccess != nil {
			if r, ok := claims.RealmAccess["roles"]; ok {
				roles = r
			}
		}
		if len(roles) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No roles found in token"})
			return
		}

		var req dto.ActionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request body",
				"errors":  []string{err.Error()},
			})
			return
		}
		c.Set("actionRequest", req)
		action := req.Action
		// Check if any role allows the action
		for _, role := range roles {
			allowedActions, exists := rolePermissions[role]
			if !exists {
				continue
			}
			for _, a := range allowedActions {
				if a == action {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied for this action"})
	}
}

// AuthMiddleware validates JWT tokens and stores user claims in the Gin context
// Performs signature, issuer, and audience checks
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		base := os.Getenv("KEYCLOAK_BASE_URL")
		realm := os.Getenv("KEYCLOAK_REALM")
		clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
		expectedIssuer := strings.TrimRight(base, "/") + "/realms/" + realm
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "MissingBearerToken"})
			return
		}
		tokenStr := strings.TrimSpace(authHeader[len("Bearer "):])
		token, err := jwt.ParseWithClaims(tokenStr, &security.KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}
			kid, _ := token.Header["kid"].(string)
			return security.GetPublicKey(ctx.Request.Context(), kid)
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "InvalidToken", "details": err.Error()})
			return
		}
		claims, _ := token.Claims.(*security.KeycloakClaims)
		// Issuer validation
		if base != "" && realm != "" && claims.Issuer != expectedIssuer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "BadIssuer", "expected": expectedIssuer, "actual": claims.Issuer})
			return
		}
		// Audience / authorized party validation:
		// Keycloak may put client id in aud OR azp depending on flow.
		strictAud := strings.ToLower(os.Getenv("KEYCLOAK_STRICT_AUD")) == "true"
		if clientID != "" {
			audOK := false
			// Audience may be a slice of strings in RegisteredClaims.Audience
			for _, a := range claims.Audience {
				if a == clientID {
					audOK = true
					break
				}
			}
			if claims.AuthorizedParty == clientID {
				audOK = true
			}
			if !audOK && strictAud {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "BadAudience", "expected": clientID, "aud": claims.Audience, "azp": claims.AuthorizedParty})
				return
			}
		}
		var roles []string
		if claims.RealmAccess != nil {
			if r, ok := claims.RealmAccess["roles"]; ok {
				roles = r
			}
		}
		logger.Info("Authenticated user:", claims.PreferredUsername, " with roles: ", claims.RealmAccess["roles"])
		ctx.Set("claims", claims)
		ctx.Set("user", claims.PreferredUsername)
		ctx.Set("roles", roles)
		if len(roles) > 0 {
			ctx.Set("role", roles[0])
		}
		reqCtx := ctx.Request.Context()
		reqCtx = context.WithValue(reqCtx, "user", claims.PreferredUsername)
		reqCtx = context.WithValue(reqCtx, "roles", roles)
		selectedRole := ""
		if len(roles) > 0 {
			for _, r := range roles {
				if r == strings.ToUpper(r) {
					selectedRole = r
					break
				}
			}
			reqCtx = context.WithValue(reqCtx, "role", selectedRole)
		}
		ctx.Request = ctx.Request.WithContext(reqCtx)

		ctx.Next()
	}
}

// MDMSRoleMiddleware checks roles based on MDMS API configuration
func MDMSRoleMiddleware(method, endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get allowed roles from pre-loaded MDMS cache
		allowedRoles := getAllowedRolesFromCache(method, endpoint)

		// If no roles found in MDMS cache, deny access
		if len(allowedRoles) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":    "No roles configured for this endpoint",
				"endpoint": endpoint,
				"method":   method,
			})
			return
		}

		// Use security.RequireRoles middleware with the cached roles
		roleMiddleware := security.RequireRoles(allowedRoles...)
		roleMiddleware(c)
	}
}

// getAllowedRolesFromCache fetches allowed roles from the pre-loaded cache
func getAllowedRolesFromCache(method, endpoint string) []string {
	uniqueIdentifier := fmt.Sprintf("%s.%s", endpoint, method)
	if roles, exists := mdmsRolePermissions[uniqueIdentifier]; exists {
		return roles
	}
	return []string{}
}

// loadAllMDMSRoles loads all MDMS roles once at startup
func loadAllMDMSRoles() (map[string][]string, error) {
	mdmsURL := config.GetConfig().MDMSURL
	url := fmt.Sprintf("%s/mdms-v2/v2?schemaCode=API_ROLES", mdmsURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers
	req.Header.Set("X-Tenant-ID", "pg")
	req.Header.Set("X-Client-Id", "test-client")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mdmsResponse MDMSResponse
	if err := json.Unmarshal(body, &mdmsResponse); err != nil {
		return nil, err
	}

	// Build cache map
	roleCache := make(map[string][]string)
	for _, mdmsData := range mdmsResponse.MDMS {
		if mdmsData.IsActive {
			roleCache[mdmsData.UniqueIdentifier] = mdmsData.Data.AllowedRoles
		}
	}

	return roleCache, nil
}

// loadRoleActionsFromMDMS loads role-action permissions from MDMS API_ROLE_ACTION schema
func loadRoleActionsFromMDMS() (map[string][]string, error) {
	mdmsURL := config.GetConfig().MDMSURL
	url := fmt.Sprintf("%s/mdms-v2/v2?schemaCode=API_ROLE_ACTION", mdmsURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers
	req.Header.Set("X-Tenant-ID", constants.DefaultJurisdiction)
	req.Header.Set("X-Client-Id", "test-client")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse as role-action data structure
	var roleActionResponse struct {
		MDMS []struct {
			ID               string         `json:"id"`
			TenantID         string         `json:"tenantId"`
			SchemaCode       string         `json:"schemaCode"`
			UniqueIdentifier string         `json:"uniqueIdentifier"`
			Data             RoleActionData `json:"data"`
			IsActive         bool           `json:"isActive"`
		} `json:"mdms"`
	}

	if err := json.Unmarshal(body, &roleActionResponse); err != nil {
		return nil, err
	}

	// Build role permissions map
	rolePermissions := make(map[string][]string)
	for _, mdmsData := range roleActionResponse.MDMS {
		if mdmsData.IsActive {
			rolePermissions[mdmsData.Data.Role] = mdmsData.Data.Actions
		}
	}

	return rolePermissions, nil
}
