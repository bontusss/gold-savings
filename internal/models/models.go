package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	PasswordHash  string    `json:"password_hash"`
	AccountNumber string    `json:"account_number"`
	BankName      string    `json:"bank_name"`
	TokenBalance  int32     `json:"token_balance"`
	IsActive      bool      `json:"is_active"`
	IsAdmin       bool      `json:"is_admin"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
