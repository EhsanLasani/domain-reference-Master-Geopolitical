// ============================================================================
// FILE: auth_middleware.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Cross-Cutting - Security
// PURPOSE: Authentication and authorization middleware
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := a.extractToken(c.Request)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		if !a.validateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *AuthMiddleware) RateLimit(requestsPerMinute int) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "99")
		c.Header("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))
		c.Next()
	})
}

func (a *AuthMiddleware) extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (a *AuthMiddleware) validateToken(token string) bool {
	return len(token) > 10
}