// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/config"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (a *Service) CreateAdminUser(ctx context.Context, email, password string) (*db.CreateAdminRow, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := a.queries.CreateAdmin(ctx, db.CreateAdminParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type UserObjectWithToken struct {
	User  *db.User
	Token string
}

func (a *Service) LoginAdmin(ctx context.Context, email, password string) (string, *db.Admin, error) {
	admin, err := a.queries.GetAdminByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token for admin
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		return "", nil, err
	}

	return signedToken, &admin, nil
}

func GenerateReferenceID(username string) string {
	// Take first 3 letters of username (or pad if less)
	trimmed := strings.ToUpper(username)
	if len(trimmed) < 3 {
		trimmed = trimmed + strings.Repeat("X", 3-len(trimmed))
	}
	letters := trimmed[:3]
	numbers := rand.Intn(9000) + 1000 // random 4-digit number
	return fmt.Sprintf("GSAV-%s-%d", letters, numbers)
}

func (a *Service) CreateUser(ctx context.Context, email, password, username, phone string) (*UserObjectWithToken, error) {
	log.Println("starting reg service")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	refID := GenerateReferenceID(username)
	fmt.Printf("user ref code is: %s", refID)

	dbUser, err := a.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Username:     username,
		Phone:        phone,
		ReferenceID:  refID,
	})
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   dbUser.ID,
		"email": dbUser.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		return nil, err
	}

	var user UserObjectWithToken
	user.User = &dbUser
	user.Token = signedToken

	return &user, nil
}

func (a *Service) DeleteUsers(ctx context.Context) error {
	err := a.queries.DeleteUsers(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *Service) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	err := a.queries.DeleteUserByID(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Service) Login(ctx context.Context, email, password string) (*UserObjectWithToken, error) {
	dbUser, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   dbUser.ID,
		"email": dbUser.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.config.JwtSecret))
	if err != nil {
		return nil, err
	}

	var user UserObjectWithToken
	user.User = &dbUser
	user.Token = signedToken

	return &user, nil
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

// ...existing code...

// SetEmailVerification sets the verification code and expiry for a user.
func (a *Service) SetEmailVerification(ctx context.Context, userID string, code string, expiry time.Time) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.queries.SetUserEmailVerification(ctx, db.SetUserEmailVerificationParams{
		ID:                    uid,
		VerificationCode:      sql.NullString{Valid: true, String: code},
		VerificationExpiresAt: sql.NullTime{Valid: true, Time: expiry},
	})
}

// VerifyEmailCode checks the code and marks the email as verified if valid and not expired.
func (a *Service) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	user, err := a.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	if user.EmailVerified {
		return false, nil // Already verified
	}
	if user.VerificationCode.String != code {
		return false, nil // Invalid code
	}
	if !user.VerificationExpiresAt.Valid || user.VerificationExpiresAt.Time.Before(time.Now()) {
		return false, nil // Expired
	}
	// Mark as verified and clear code
	err = a.queries.MarkUserEmailVerified(ctx, db.MarkUserEmailVerifiedParams{
		ID:               user.ID,
		EmailVerified:    true,
		VerificationCode: sql.NullString{String: "", Valid: false},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// ...existing code...
