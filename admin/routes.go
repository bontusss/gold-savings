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
	"gold-savings/internal/services"

	"github.com/gin-gonic/gin"
)

// routes TODO:

func SetupRoutes(router *gin.RouterGroup, authService *auth.Service, queries *db.Queries, admin *services.Admin) {
	authHandler := handlers.NewAuthHandler(authService)
	dashboardHandler := handlers.NewDashboardHandler(authService, queries, admin)

	// Auth routes
	router.GET("/login", authHandler.ShowLogin)
	router.POST("/login", authHandler.Login)
	router.GET("/register", authHandler.ShowRegister)
	router.POST("/register", authHandler.Register)
	router.POST("/logout", authHandler.Logout)
	// delete after test
	router.POST("/api/plan", dashboardHandler.CreateInvestmentPlan)
	router.GET("/api/approve-payment/:id", dashboardHandler.ApprovePayment)
	router.GET("/api/decline-payment/:id", dashboardHandler.DeclinePayment)
	router.GET("/api/approve-withdraw/:id", dashboardHandler.ApproveWithdrawal)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/dashboard", dashboardHandler.ShowDashboard)
		protected.GET("/api/data", dashboardHandler.GetData)
		protected.GET("/api/users", dashboardHandler.ListUsers)
		protected.GET("/api/plan", dashboardHandler.ShowCreatePlan)

		// protected.POST("/api/plan", dashboardHandler.CreateInvestmentPlan)
	}
}
