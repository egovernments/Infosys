// Package security provides JWT and role-based security utilities for the API
package security

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// keycloakCerts holds the JWKS keys fetched from Keycloak for JWT validation
type keycloakCerts struct {
	Keys []struct {
		KeyID            string   `json:"kid"`
		KeyType          string   `json:"kty"`
		Algorithm        string   `json:"alg"`
		PublicKeyUse     string   `json:"use"`
		Modulus          string   `json:"n"`
		Exponent         string   `json:"e"`
		CertificateChain []string `json:"x5c"`
	} `json:"keys"`
}

var (
	jwksCache       keycloakCerts      // Cached JWKS keys
	jwksFetchedAt   time.Time          // Last fetch time
	jwksMutex       sync.RWMutex       // Mutex for thread safety
	jwksTTL         = 10 * time.Minute // Cache TTL
	ErrUnauthorized = errors.New("unauthorized")
)

// fetchJWKS retrieves the JSON Web Key Set (JWKS) from the Keycloak realm endpoint
func fetchJWKS(ctx context.Context) (keycloakCerts, error) {
	base := os.Getenv("KEYCLOAK_BASE_URL")
	realm := os.Getenv("KEYCLOAK_REALM")
	if base == "" || realm == "" {
		return keycloakCerts{}, fmt.Errorf("KEYCLOAK_BASE_URL or KEYCLOAK_REALM not set")
	}
	url := strings.TrimRight(base, "/") + "/realms/" + realm + "/protocol/openid-connect/certs"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return keycloakCerts{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return keycloakCerts{}, fmt.Errorf("failed to fetch JWKS: %s", resp.Status)
	}
	var kc keycloakCerts
	if err := json.NewDecoder(resp.Body).Decode(&kc); err != nil {
		return keycloakCerts{}, err
	}
	return kc, nil
}

// getJWKS retrieves the JWKS with in-memory caching to avoid repeated HTTP requests
func getJWKS(ctx context.Context) (keycloakCerts, error) {
	jwksMutex.RLock()
	if time.Since(jwksFetchedAt) < jwksTTL && len(jwksCache.Keys) > 0 {
		defer jwksMutex.RUnlock()
		return jwksCache, nil
	}
	jwksMutex.RUnlock()

	jwksMutex.Lock()
	defer jwksMutex.Unlock()
	// Double-check condition after acquiring write lock
	if time.Since(jwksFetchedAt) < jwksTTL && len(jwksCache.Keys) > 0 {
		return jwksCache, nil
	}
	kc, err := fetchJWKS(ctx)
	if err != nil {
		return keycloakCerts{}, err
	}
	jwksCache = kc
	jwksFetchedAt = time.Now()
	return kc, nil
}

// parseRSAPublicKey constructs an RSA public key from base64url-encoded modulus and exponent values
func parseRSAPublicKey(nB64, eB64 string) (*rsa.PublicKey, error) {
	// n and e are base64url encoded per RFC 7517
	if nB64 == "" || eB64 == "" {
		return nil, errors.New("missing modulus or exponent")
	}
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode n: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode e: %w", err)
	}
	// Exponent is big-endian integer; typically small (e.g., 65537)
	var eInt int
	switch len(eBytes) {
	case 3: // common for 65537 (0x01 0x00 0x01)
		eInt = int(binary.BigEndian.Uint32(append([]byte{0x00}, eBytes...)))
	case 4:
		eInt = int(binary.BigEndian.Uint32(eBytes))
	case 1:
		eInt = int(eBytes[0])
	default:
		// generic parse
		eInt = 0
		for _, b := range eBytes {
			eInt = eInt<<8 + int(b)
		}
	}
	modulus := new(big.Int).SetBytes(nBytes)
	return &rsa.PublicKey{N: modulus, E: eInt}, nil
}

// GetPublicKey retrieves the RSA public key for a given key ID (kid) from the JWKS
func GetPublicKey(ctx context.Context, keyId string) (interface{}, error) {
	jwks, err := getJWKS(ctx)
	if err != nil {
		return nil, err
	}
	for _, key := range jwks.Keys {
		if key.KeyID != keyId {
			continue
		}
		// Prefer certificate chain if present
		if len(key.CertificateChain) > 0 {
			pem := "-----BEGIN CERTIFICATE-----\n" + key.CertificateChain[0] + "\n-----END CERTIFICATE-----\n"
			pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(extractPublicKeyPEM(pem)))
			if err == nil {
				return pubKey, nil
			}
			// fall through to n/e if cert parsing fails
		}
		if key.Modulus != "" && key.Exponent != "" {
			rsaKey, err := parseRSAPublicKey(key.Modulus, key.Exponent)
			if err == nil {
				return rsaKey, nil
			}
			return nil, err
		}
		return nil, errors.New("jwks entry missing x5c and n/e")
	}
	return nil, fmt.Errorf("public key not found for kid %s", keyId)
}

// extractPublicKeyPEM returns the certificate PEM as-is (for compatibility)
func extractPublicKeyPEM(certPEM string) string { return certPEM }

// KeycloakClaims represents the JWT claims structure returned by Keycloak
type KeycloakClaims struct {
	jwt.RegisteredClaims
	PreferredUsername string              `json:"preferred_username"`
	Email             string              `json:"email"`
	RealmAccess       map[string][]string `json:"realm_access"`
	ResourceAccess    map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	AuthorizedParty string `json:"azp"` // client id when acting on behalf of
}

// RequireRoles returns a Gin middleware that checks if the authenticated user has any of the required roles
func RequireRoles(required ...string) gin.HandlerFunc {
	reqSet := make(map[string]struct{}, len(required))
	for _, r := range required {
		reqSet[r] = struct{}{}
	}
	return func(ctx *gin.Context) {
		v, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "AuthClaimsMissing"})
			return
		}
		claims := v.(*KeycloakClaims)
		roles := collectRoles(claims)
		for r := range roles {
			if _, ok := reqSet[r]; ok {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden", "requiredRoles": required})
	}
}

// collectRoles extracts all roles (realm and client roles) from the Keycloak JWT claims
func collectRoles(claims *KeycloakClaims) map[string]struct{} {
	out := make(map[string]struct{})
	if claims == nil {
		return out
	}
	// realm roles
	if claims.RealmAccess != nil {
		for _, r := range claims.RealmAccess["roles"] {
			out[r] = struct{}{}
		}
	}
	// client roles
	for _, ra := range claims.ResourceAccess {
		for _, r := range ra.Roles {
			out[r] = struct{}{}
		}
	}
	return out
}

// GetUserID extracts the user ID (preferred username) from the authenticated user's JWT claims
func GetUserID(ctx *gin.Context) string {
	v, ok := ctx.Get("claims")
	if !ok {
		return ""
	}
	cl := v.(*KeycloakClaims)
	return cl.PreferredUsername
}
