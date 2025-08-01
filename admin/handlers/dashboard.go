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
	"strconv"

	"github.com/gin-gonic/gin"
)

// handler.go TODO:

type DashboardHandler struct {
	authService *auth.Service
	queries     *db.Queries
	admin       *services.Admin
}

var serverError = "an error occurred, try again"

func NewDashboardHandler(authService *auth.Service, queries *db.Queries, admin *services.Admin) *DashboardHandler {
	return &DashboardHandler{authService: authService, queries: queries, admin: admin}
}

func (h *DashboardHandler) ShowDashboard(c *gin.Context) {

	activeUsersCount, err := h.queries.CountActiveUsers(c)
	if err != nil {
		log.Printf("error getting active users count: %v", err)
		c.String(500, "error fetching active user count")
		return
	}

	ts, err := h.queries.GetTotalSavingsInApp(c)
	if err != nil {
		log.Printf("error getting total savings: %v", err)
		c.String(500, "error fetching total savings")
		return
	}

	ia, err := h.queries.GetTotalInvestmentInApp(c)
	if err != nil {
		log.Printf("error getting total savings: %v", err)
		c.String(500, "error fetching total savings")
		return
	}

	tt, err := h.queries.GetTotalTokens(c)
	if err != nil {
		log.Printf("error getting total savings: %v", err)
		c.String(500, "error fetching total savings")
		return
	}
	requests, _ := h.queries.ListAllPayoutRequests(c)
	topUsers, _ := h.queries.ListUsersByTotalSavingsDesc(c)
	trans, _ := h.queries.ListPendingDepositTransactionsWithUser(c)
	plans, _ := h.queries.GetAllInvestmentPlans(c)
	i, _ := h.queries.ListInvestmentsWithUserAndPlan(c)
	deposits, _ := h.queries.ListTransactionsByCategory(c, "savings")
	wdw, _ := h.queries.ListPendingWithdrawalTransactionsWithUser(c)
		
	err = components.DashboardT(activeUsersCount, ts, ia, tt, topUsers, plans, requests, trans, i, deposits, wdw).Render(c, c.Writer)
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

func (h *DashboardHandler) ShowCreatePlan(c *gin.Context) {
	err := components.CreatePlan().Render(c, c.Writer)
	if err != nil {
		log.Printf("Error rendering create plan template: %v", err)
		c.String(http.StatusInternalServerError, "Error rendering create plans page")
		return
	}
}

func (h *DashboardHandler) CreateInvestmentPlan(c *gin.Context) {
	var req db.CreateInvestmentPlanParams
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("create plan request error: %v", err)
		c.JSON(400, serverError)
		return
	}
	_, err := h.admin.CreateInvestmentPlan(c, req.Name, req.InterestRate, req.MinAmount, req.MaxAmount)
	if err != nil {
		log.Printf("error creating plan: %v", err)
		c.JSON(500, serverError)
		return
	}
	c.JSON(200, "plan created")
}

func (h *DashboardHandler) ApprovePayment(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		log.Printf("invalid transaction ID: %v", err)
		c.JSON(400, "invalid transaction ID")
		return
	}

	transaction, err := h.queries.GetTransactionByID(c, int32(transID))
	if err != nil {
		log.Printf("error getting transaction: %v", err)
		c.JSON(400, serverError)
		return
	}

	switch transaction.Category {
	case "savings":
		err = h.admin.ApprovePayment(c, int32(transID), "approved", "Approved by admin via email link")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	case "investment":
		err = h.admin.ApproveInvestment(c, int32(transID), "approved", "Approved by admin")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	case "token":
		err = h.admin.ApproveTokenRequest(c, int32(transID), "approved", "Approved by admin")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	default:
		log.Println("Invalid transaction category")
		c.JSON(500, serverError)
		return
	}

	// Approve directly, no JSON body needed

	c.JSON(200, "transaction status updated")
}

func (h *DashboardHandler) DeclinePayment(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		log.Printf("invalid transaction ID: %v", err)
		c.JSON(400, "invalid transaction ID")
		return
	}

	transaction, err := h.queries.GetTransactionByID(c, int32(transID))
	if err != nil {
		log.Printf("error getting transaction: %v", err)
		c.JSON(400, serverError)
		return
	}

	switch transaction.Category {
	case "savings":
		err = h.admin.DeclinePayment(c, int32(transID), "declined", "Declined by admin via email link")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	case "investment":
		err = h.admin.DeclineInvestment(c, int32(transID), "declined", "Declined by admin via email link")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	}

	c.JSON(200, "transaction status updated")
}

func (h *DashboardHandler) ApproveWithdrawal(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		log.Printf("invalid transaction ID: %v", err)
		c.JSON(400, "invalid transaction ID")
		return
	}

	transaction, err := h.queries.GetTransactionByID(c, int32(transID))
	if err != nil {
		log.Printf("error getting transaction: %v", err)
		c.JSON(400, serverError)
		return
	}

	switch transaction.Category {
	case "savings":
		err = h.admin.ApproveWithdrawal(c, int32(transID), "approved", "Approved by admin via email link")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	case "investment":
		err = h.admin.ApproveInvestmentWithdraw(c, int32(transID), "approved", "Approved by admin")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	case "token":
		err = h.admin.ApproveTokenRequest(c, int32(transID), "approved", "Approved by admin")
		if err != nil {
			log.Printf("error updating transaction status: %v", err)
			c.JSON(500, serverError)
			return
		}
	default:
		log.Println("Invalid transaction category")
		c.JSON(500, serverError)
		return
	}

	// Approve directly, no JSON body needed

	c.JSON(200, "transaction status updated")
}

func isHTMX(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}
