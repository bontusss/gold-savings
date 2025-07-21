// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/8/2025

package middleware

import (
	"gold-savings/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// auth TODO:

func JWTAuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Check if the header has the Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be Bearer {token}",
			})
			return
		}

		token := parts[1]

		// Validate the token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		subFloat, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			return
		}

		emailStr, ok := claims["email"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email in token",
			})
			return
		}

		isAdmin, _ := claims["is_admin"].(bool)

		c.Set("userID", int32(subFloat))
		c.Set("email", emailStr)
		c.Set("isAdmin", isAdmin)

		c.Next()
	}
}
