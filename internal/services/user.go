package services

import (
	"context"
	"database/sql"
	db "gold-savings/db/sqlc"

	"github.com/google/uuid"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}
func (s *UserService) GetUserByEmail(email string) (*db.GetUserByEmailRow, error) {
	user, err := s.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) ListUsers(id string) (*[]db.User, error) {
	users, err := s.queries.ListUsers(context.Background())
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (s *UserService) UpdateStatus(userid string, isActive bool) error {
	err := s.queries.UpdateUserStatus(context.Background(), db.UpdateUserStatusParams{
		ID:       uuid.MustParse(userid),
		IsActive: sql.NullBool{Bool: isActive, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}
