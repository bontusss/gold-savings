// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package admin

import (
	"gold-savings/admin/handlers"
	"gold-savings/admin/middleware"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/auth"

	"github.com/gin-gonic/gin"
)

// routes TODO:

func SetupRoutes(router *gin.RouterGroup, authService *auth.Service, queries *db.Queries) {
	authHandler := handlers.NewAuthHandler(authService)
	dashboardHandler := handlers.NewDashboardHandler(authService, queries)

	// Auth routes
	router.GET("/login", authHandler.ShowLogin)
	router.POST("/login", authHandler.Login)
	router.GET("/register", authHandler.ShowRegister)
	router.POST("/register", authHandler.Register)
	router.POST("/logout", authHandler.Logout)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/dashboard", dashboardHandler.ShowDashboard)
		protected.GET("/api/data", dashboardHandler.GetData)
		protected.GET("/api/users", dashboardHandler.ListUsers)
	}
}
