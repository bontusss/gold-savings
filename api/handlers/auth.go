// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/8/2025

package handlers

import (
	"gold-savings/api/utils"
	"gold-savings/internal/auth"
	"gold-savings/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// auth TODO:

type AuthHandler struct {
	authService *auth.Service
	config      *config.Config
}

var serverError = "an error ocurred, try again"

func NewAuthHandler(a *auth.Service, config *config.Config) *AuthHandler {
	return &AuthHandler{authService: a, config: config}
}

// Login handles API login requests (for mobile app)
func (h *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.authService.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		log.Printf("error logging in: %v", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error":  "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   user,
	})
}

// Register handles API registration requests (for mobile app)
func (h *AuthHandler) Register(c *gin.Context) {
	log.Println("starting registrattion")
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Username string `json:"username" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("inavalid register data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "an error occured, try again",
		})
		return
	}

	user, err := h.authService.CreateUser(c.Request.Context(), request.Email, request.Password, request.Username, request.Phone)
	if err != nil {
		log.Printf("error creating user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Registration failed: user already exists",
		})
		return
	}

	// Generate verification code and expiry
	code := utils.GenerateOTP()
	expiry := time.Now().Add(10 * time.Minute)

	// Save code and expiry to user (implement this in your service/repo)
	err = h.authService.SetEmailVerification(c.Request.Context(), user.User.ID.String(), code, expiry)
	if err != nil {
		log.Printf("error saving verification code: %v", err)
	}

	// Send verification email
	emailBody, _ := utils.RenderEmailTemplate("templates/verify_email.html", map[string]any{
		"Username": user.User.Username,
		"Code":     code,
	})
	plunk := utils.Plunk{HttpClient: http.DefaultClient, Config: h.config}
	err = plunk.SendEmail(user.User.Email, "Verify your Gold Savings account", emailBody)
	if err != nil {
		log.Printf("error sending verification email: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"user":   user,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ok, err := h.authService.VerifyEmailCode(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid or expired code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}

func (h AuthHandler) DeleteUsers(c *gin.Context) {
	err := h.authService.DeleteUsers(c)
	if err != nil {
		log.Printf("error deleteing user: %v", err)
		c.JSON(500, serverError)
	}
}

// func (h AuthHandler) DeleteUserByID(c *gin.Context) {

// 	err := h.authService.DeleteUserByID(c)
// 	if err != nil {
// 		log.Printf("error deleteing user: %v", err)
// 		c.JSON(500, serverError)
// 	}
// }
