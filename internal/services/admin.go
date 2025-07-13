package services

import (
	"context"
	db "gold-savings/db/sqlc"
)

type Admin struct {
	queries *db.Queries
}

func NewAdminService(q *db.Queries) *Admin {
	return &Admin{q}
}

func (a *Admin) CreateInvestmentPlan(ctx context.Context, name, interestRate, min_amount, max_amount string) (*db.InvestmentPlan, error) {
	args := db.CreateInvestmentPlanParams{
		Name:         name,
		InterestRate: interestRate,
		MinAmount:    min_amount,
		MaxAmount:    max_amount,
	}
	plan, err := a.queries.CreateInvestmentPlan(ctx, args)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}
