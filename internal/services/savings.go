// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/13/2025

package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	db "gold-savings/db/sqlc"
	"math/rand"
	"time"
)

// savings TODO:

type Savings struct {
	queries *db.Queries
}

type SavingsPlanStatus string

const (
	SavingsPlanActive    SavingsPlanStatus = "active"
	SavingsPlanCancelled SavingsPlanStatus = "cancelled"
	SavingsPlanCompleted SavingsPlanStatus = "completed"
)

func (s *Savings) GeneratePlanRef(ctx context.Context, userId uuid.UUID) (string, error) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)
	user, err := s.queries.GetUser(ctx, userId)
	if err != nil {
		return "", err
	}
	planRef := fmt.Sprintf("GSAV-%s-%04d", user.FirstName, randomNumber)
	return planRef, nil
}

func (s *Savings) CreatePlan(ctx context.Context, userID uuid.UUID, target, savingsAmountPerFrequency decimal.Decimal, duration int32, frequency string) (*db.SavingsPlan, error) {
	ref, err := s.GeneratePlanRef(ctx, userID)
	if err != nil {
		return nil, err
	}
	params := db.CreateSavingsPlanParams{
		UserID:           userID,
		PlanRef:          ref,
		TargetAmount:     target.String(),
		CurrentAmount:    "0",
		DurationDays:     duration,
		SavingsFrequency: frequency,
		SavingsAmount:    savingsAmountPerFrequency.String(),
		Status:           string(SavingsPlanActive),
	}
	plan, err := s.queries.CreateSavingsPlan(ctx, params)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *Savings) GetSavingsPlan(ctx context.Context, id uuid.UUID) (*db.SavingsPlan, error) {
	plan, err := s.queries.GetSavingsPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *Savings) GetSavingsPlanByRef(ctx context.Context, ref string) (*db.SavingsPlan, error) {
	plan, err := s.queries.GetSavingsPlanByPlanRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *Savings) GetUserPlans(ctx context.Context, userID uuid.UUID) (*[]db.SavingsPlan, error) {
	plans, err := s.queries.ListSavingsPlansByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &plans, nil
}

func (s *Savings) UpdatePlanAmount(ctx context.Context, planID uuid.UUID, amount decimal.Decimal) (*db.SavingsPlan, error) {
	params := db.UpdateSavingsPlanAmountParams{
		ID:            planID,
		CurrentAmount: amount.String(),
	}
	plan, err := s.queries.UpdateSavingsPlanAmount(ctx, params)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *Savings) UpdateStatus(ctx context.Context, planID uuid.UUID, status string) (*db.SavingsPlan, error) {
	params := db.UpdateSavingsPlanStatusParams{
		ID:     planID,
		Status: status,
	}
	plan, err := s.queries.UpdateSavingsPlanStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *Savings) ListActivePlans(ctx context.Context) (*[]db.SavingsPlan, error) {
	plans, err := s.queries.ListActiveSavingsPlans(ctx)
	if err != nil {
		return nil, err
	}
	return &plans, nil
}

func (s *Savings) DeletePlan(ctx context.Context, planID uuid.UUID) error {
	err := s.queries.DeleteSavingsPlan(ctx, planID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Savings) ListAllSavingsPlans(ctx context.Context) (*[]db.SavingsPlan, error) {
	plans, err := s.queries.ListAllSavingsPlans(ctx)
	if err != nil {
		return nil, err
	}
	return &plans, nil
}
