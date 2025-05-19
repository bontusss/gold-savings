// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/8/2025

package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// data TODO:

type DataHandler struct{}

func NewDataHandler() *DataHandler {
	return &DataHandler{}
}

// HealthCheck is a simple endpoint to verify the API is running
func (h *DataHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"message": "API is healthy",
	})
}

// GetProtectedData is an example of a protected endpoint
func (h *DataHandler) GetProtectedData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User information not found",
		})
		return
	}

	email, _ := c.Get("email")
	isAdmin, _ := c.Get("isAdmin")

	c.JSON(http.StatusOK, gin.H{
		"message": "This is protected data",
		"user": gin.H{
			"id":      userID,
			"email":   email,
			"isAdmin": isAdmin,
		},
		"data": []gin.H{
			{"id": 1, "name": "Item 1"},
			{"id": 2, "name": "Item 2"},
			{"id": 3, "name": "Item 3"},
		},
	})
}
