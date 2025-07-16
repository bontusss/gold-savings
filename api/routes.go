// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package api

import (
	"gold-savings/api/handlers"
	"gold-savings/api/middleware"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/auth"
	"gold-savings/internal/config"
	"gold-savings/internal/services"

	"github.com/gin-gonic/gin"
)

// routes TODO:

func SetupRoutes(router *gin.RouterGroup, authService *auth.Service, queries *db.Queries, config *config.Config, userService *services.UserService) {
	authHandler := handlers.NewAuthHandler(authService, config)
	userHandler := handlers.NewUserHandler(userService)
	// Public routes
	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)
	router.POST("/verify_email", authHandler.VerifyEmail)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(authService))
	{
		protected.GET("/plans", userHandler.GetAllInvestmentPlans)
		protected.POST("/savings/payment_request", userHandler.CreateSavingsPaymentRequest)
		protected.POST("/investment/payment_request", userHandler.CreateInvestmentPaymentRequest)
		protected.GET("/transactions/savings", userHandler.ListUserSavingsTransactions)
		protected.GET("/transactions/investment", userHandler.ListUserInvestmentTransactions)
		protected.POST("/investment", userHandler.CreateInvestment)
		protected.GET("/investments", userHandler.ListUserInvestments)
		protected.GET("/user/me", userHandler.GetUser)
	}
}
