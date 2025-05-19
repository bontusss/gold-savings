// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/8/2025

package handlers

import (
	"github.com/gin-gonic/gin"
	"gold-savings/internal/auth"
	"net/http"
)

// auth TODO:

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(a *auth.Service) *AuthHandler {
	return &AuthHandler{authService: a}
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

	token, err := h.authService.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Register handles API registration requests (for mobile app)
func (h *AuthHandler) Register(c *gin.Context) {
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

	user, err := h.authService.CreateAdminUser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Registration failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}
