package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rezbow/ecommerce/internal/platform/authentication"
	"github.com/rezbow/ecommerce/internal/platform/config"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. EXTRACT: Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// 2. PARSE: Check for "Bearer " prefix and extract the token string
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format (Expected: Bearer <token>)"})
			return
		}
		tokenString := parts[1]

		// 3. VALIDATE: Call the platform JWT utility to parse and validate
		claims, err := authentication.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			// This catches expired, invalid signature, or malformed tokens
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			return
		}

		// 4. INJECT: Add user data from claims into Gin's context
		// Use the custom context keys for safe retrieval later
		c.Set("userId", claims.UserId)
		c.Set("isAdmin", claims.IsAdmin)

		// 5. CONTINUE: Token is valid, proceed to the next handler
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdminVal, exists := c.Get("isAdmin")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "unauthorized"})
			return
		}

		isAdmin, ok := isAdminVal.(bool)
		if !ok || !isAdmin {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Requires Administrator privileges"})
			return
		}
		c.Next()
	}
}
