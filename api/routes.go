// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package api

import (
	"github.com/gin-gonic/gin"
	"gold-savings/api/handlers"
	"gold-savings/api/middleware"
	"gold-savings/internal/auth"
)

// routes TODO:

func SetupRoutes(router *gin.RouterGroup, authService *auth.Service) {
	authHandler := handlers.NewAuthHandler(authService)
	dataHandler := handlers.NewDataHandler()
	// Public routes
	router.GET("/health", dataHandler.HealthCheck)
	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(authService))
	{
		protected.GET("/data", dataHandler.GetProtectedData)
	}
}
