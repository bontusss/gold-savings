// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package auth

import (
	"context"
	"database/sql"
	"errors"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// auth TODO:

type Service struct {
	queries *db.Queries
	config  *config.Config
}

func NewAuthService(queries *db.Queries, config *config.Config) *Service {
	return &Service{
		queries: queries,
		config:  config,
	}
}

func (a *Service) CreateAdminUser(ctx context.Context, email, password string) (*db.CreateAdminUserRow, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := a.queries.CreateAdminUser(ctx, db.CreateAdminUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsAdmin:      sql.NullBool{Valid: true, Bool: true},
		FirstName:    "admin",
		LastName:     "admin",
		Phone:        "+12405551212",
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Ensure is_admin is a plain boolean
	isAdmin := user.IsAdmin.Valid && user.IsAdmin.Bool

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"email":    user.Email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(a.config.JwtSecret))
}

func (a *Service) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.config.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
