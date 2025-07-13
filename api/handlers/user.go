package handlers

import (
	"gold-savings/api/utils"
	"gold-savings/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type userHandler struct {
	userService *services.UserService
}

func NewUserHandler(u *services.UserService) *userHandler {
	return &userHandler{u}
}

func (s *userHandler) GetAllInvestmentPlans(c *gin.Context) {
	plans, err := s.userService.GetInvestmentPlans(c)
	if err != nil {
		log.Printf("error getting investment plans: %v", err)
		c.JSON(500, serverError)
		return
	}

	c.JSON(200, plans)
}

func (s *userHandler) CreateSavingsPaymentRequest(c *gin.Context) {
	var req struct {
		Amount        string `json:"amount" binding:"required"`
		BankName      string `json:"bank_name" binding:"required"`
		AccountName   string `json:"account_name" binding:"required"`
		AccountNumber string `json:"account_number" binding:"required"`
	}

	userID, err := utils.GetActiveUserID(c)
	if err != nil {
		log.Printf("error getting active user ID: %v", err)
		c.JSON(500, serverError)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	amountDecimal, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Printf("error converting amount to decimal: %v", err)
		c.JSON(400, gin.H{"error": "invalid amount format"})
		return
	}

	payment, err := s.userService.CreateSavingsPaymentRequest(c, userID, amountDecimal, req.BankName, req.AccountName, req.AccountNumber)
	if err != nil {
		log.Printf("error creating payment request: %v", err)
		c.JSON(500, serverError)
		return
	}

	transaction, err := s.userService.CreateTransaction(c, userID, amountDecimal, payment.Type)
	if err != nil {
		log.Printf("error creating transaction: %v", err)
		c.JSON(500, serverError)
		return
	}

	c.JSON(201, transaction)
}

func (s *userHandler) CreateInvestmentPaymentRequest(c *gin.Context) {
	var req struct {
		RefCode     string `json:"ref_code" binding:"required"`
		Amount      string `json:"amount" binding:"required"`
		BankName    string `json:"bank_name" binding:"required"`
		AccountName string `json:"account_name" binding:"required"`
	}

	userID, err := utils.GetActiveUserID(c)
	if err != nil {
		log.Printf("error getting active user ID: %v", err)
		c.JSON(500, serverError)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// check if refcode for that investment exists
	iID, err := s.userService.GetInvestmentByRefCode(c, req.RefCode)
	if err != nil {
		log.Printf("GetPayoutRequestByRefCode error: %v", err)
		c.JSON(400, gin.H{"error": "invalid ref_code"})
		return
	}

	amountDecimal, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Printf("error converting amount to decimal: %v", err)
		c.JSON(400, gin.H{"error": "invalid amount format"})
		return
	}

	// Convert uuid.UUID to uuid.NullUUID
	investmentID := uuid.NullUUID{UUID: iID.ID, Valid: true}
	payment, err := s.userService.CreateInvestmentPaymentRequest(c, userID, investmentID, amountDecimal, req.BankName, req.AccountName)
	if err != nil {
		log.Printf("error creating payment request: %v", err)
		c.JSON(500, serverError)
		return
	}

	c.JSON(201, payment)
}
func (s *userHandler) ListUserSavingsTransactions(c *gin.Context) {
	userID, err := utils.GetActiveUserID(c)
	if err != nil {
		log.Printf("error getting active user ID: %v", err)
		c.JSON(500, serverError)
		return
	}

	transactions, err := s.userService.GetUserSavingsTransactions(c, userID)
	if err != nil {
		log.Printf("error getting user savings transactions: %v", err)
		c.JSON(500, serverError)
		return
	}

	c.JSON(200, transactions)
}

func (s *userHandler) ListUserInvestmentTransactions(c *gin.Context) {
	userID, err := utils.GetActiveUserID(c)
	if err != nil {
		log.Printf("error getting active user ID: %v", err)
		c.JSON(500, serverError)
		return
	}

	transactions, err := s.userService.GetUserInvestmentTransactions(c, userID)
	if err != nil {
		log.Printf("error getting user investment transactions: %v", err)
		c.JSON(500, serverError)
		return
	}

	c.JSON(200, transactions)
}

func (s *userHandler) CreateInvestment(c *gin.Context) {
	var req struct {
		PlanID int32  `json:"plan_id" binding:"required"`
		Amount string `json:"amount" binding:"required"`
	}

	userID, err := utils.GetActiveUserID(c)
	if err != nil {
		log.Printf("error getting active user ID: %v", err)
		c.JSON(500, serverError)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	amountDecimal, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Printf("error converting amount to decimal: %v", err)
		c.JSON(400, gin.H{"error": "invalid amount format"})
		return
	}

	investment, err := s.userService.CreateInvestment(c, userID, req.PlanID, amountDecimal)
	if err != nil {
		log.Printf("error creating investment: %v", err)
		c.JSON(500, err)
		return
	}

	c.JSON(201, investment)
}

func (s *userHandler) ListUserInvestments(c *gin.Context) {
	userID, err := utils.GetActiveUserID(c)
	if err != nil {
		log.Printf("error getting active user ID: %v", err)
		c.JSON(500, serverError)
		return
	}

	investments, err := s.userService.GetUserInvestmentsWithPlan(c, userID)
	if err != nil {
		log.Printf("error getting user investments: %v", err)
		c.JSON(500, serverError)
		return
	}

	c.JSON(200, investments)
}
