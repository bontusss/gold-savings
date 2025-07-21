package utils

import (
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

type LoggedInUser struct {
	ID      int32
	Email   string
	IsAdmin bool
}

type TokenObject struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"user_email"`
	Verified bool   `json:"user_verified"`
}

func GetLoggedInUser(c *gin.Context) (*LoggedInUser, error) {
	id, exists := c.Get("userID")
	if !exists {
		return nil, fmt.Errorf("error: not authorized to access this resource")
	}
	email, _ := c.Get("email")

	return &LoggedInUser{
		ID:    id.(int32),
		Email: email.(string),
	}, nil
}

func GetActiveUserID(c *gin.Context) (userID int32, err error) {
	user, err := GetLoggedInUser(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomLetters(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomNumbers(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += fmt.Sprintf("%d", rand.Intn(10))
	}
	return s
}

func GenerateCode() string {
	return fmt.Sprintf("GLD-%s-%s", randomLetters(4), randomNumbers(5))
}
