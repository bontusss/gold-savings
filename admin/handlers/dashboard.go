// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package handlers

import (
	"gold-savings/admin/components"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/auth"
	"gold-savings/internal/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handler.go TODO:

type DashboardHandler struct {
	authService *auth.Service
	queries     *db.Queries
	savings     *services.Savings
}

func NewDashboardHandler(authService *auth.Service, queries *db.Queries) *DashboardHandler {
	return &DashboardHandler{authService: authService, queries: queries}
}

func (h *DashboardHandler) ShowDashboard(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		// Log the error for debugging
		c.String(http.StatusUnauthorized, "Token validation failed: %v", err)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	// Example usage of claims (e.g., logging user email)
	if email, ok := claims["email"].(string); ok {
		c.Set("user_email", email)
	}

	activeUsersCount, err := h.queries.CountActiveNonAdminUsers(c)
	if err != nil {
		log.Printf("error getting active users count: %v", err)
		c.String(500, "error fetching active user count")
		return
	}
	totalAmountInAllActiveSavings, err := h.queries.SumActiveSavingsPlans(c)
	if err != nil {
		log.Printf("error getting active users count: %v", err)
		c.String(500, "error fetching active user count")
		return
	}

	err = components.DashboardT(activeUsersCount, totalAmountInAllActiveSavings).Render(c, c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering dashboard page")
		return
	}
}

func (h *DashboardHandler) GetData(c *gin.Context) {
	// This is a protected endpoint that makes a GET request
	data := map[string]any{
		"message": "This is protected data",
		"status":  "success",
	}

	c.JSON(http.StatusOK, data)
}

func (h *DashboardHandler) ListUsers(c *gin.Context) {
	users, err := h.queries.ListUsers(c)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		c.String(http.StatusInternalServerError, "Error fetching users")
		return
	}

	// Render the Users template with the fetched users
	err = components.ListUsers("User List", users).Render(c, c.Writer)
	if err != nil {
		log.Printf("Error rendering users template: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering users page")
		return
	}
}

func (h *DashboardHandler) ListSavingPlans(c *gin.Context) {
	_, err := h.savings.ListAllSavingsPlans(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching savings plans")
		return
	}

	//	todo: create and render the list in dashboard
}

//func (h *DashboardHandler) GetActiveUsersCount(c *gin.Context) {
//	activeUsersCount, err := h.queries.CountActiveNonAdminUsers(c)
//	if err != nil {
//		log.Printf("error getting active users count: %v", err)
//		c.String(500, "error fetching active user count")
//		return
//	}
//
//	err = components.DashboardT(activeUsersCount).Render(c, c.Writer)
//	if err != nil {
//		c.String(http.StatusInternalServerError, "Error rendering dashboard")
//		return
//	}
//}
