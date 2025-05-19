package middleware

import (
	"log"
	"net/http"

	"gold-savings/internal/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware protects routes by validating the auth_token cookie
func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth_token cookie
		token, err := c.Cookie("auth_token")
		if err != nil {
			log.Printf("No auth_token cookie: %v", err)
			c.Header("HX-Redirect", "/admin/login")
			// Fallback to standard redirect if HX-Redirect isn't handled
			c.Redirect(http.StatusFound, "/admin/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			c.Header("HX-Redirect", "/admin/login")
			c.Redirect(http.StatusFound, "/admin/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check is_admin claim
		isAdmin, ok := claims["is_admin"].(bool)
		if !ok || !isAdmin {
			log.Printf("User is not admin: ok=%v, isAdmin=%v", ok, isAdmin)
			c.Header("HX-Redirect", "/admin/login")
			c.Redirect(http.StatusFound, "/admin/login")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set user context
		c.Set("userID", claims["sub"])
		c.Set("email", claims["email"])
		c.Next()
	}
}
