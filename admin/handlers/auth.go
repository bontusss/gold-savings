// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package handlers

import (
	"gold-savings/admin/components"
	"gold-savings/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// auth TODO:

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(a *auth.Service) *AuthHandler {
	return &AuthHandler{authService: a}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	err := components.Login("").Render(c, c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering login page")
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, _, err := h.authService.LoginAdmin(c.Request.Context(), email, password)
	if err != nil {
		err = components.Login("invalid credentials").Render(c, c.Writer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error rendering login page")
		}
		return
	}

	c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/admin/dashboard")
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusFound, "/admin/dashboard")
	}
}

func (h *AuthHandler) ShowRegister(c *gin.Context) {
	err := components.Register("").Render(c, c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering login page")
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	_, err := h.authService.CreateAdminUser(c.Request.Context(), email, password)
	if err != nil {
		err := components.Register(err.Error()).Render(c, c.Writer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error rendering Register page")
			return
		}
		return
	}

	c.Header("HX-Redirect", "/admin/login")
	c.Redirect(http.StatusFound, "/admin/login") // Ensure consistent redirection to login
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Header("HX-Redirect", "/admin/login")
	c.Status(http.StatusOK)
}
