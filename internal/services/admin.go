package services

import (
	"context"
	"database/sql"
	"errors"
	"gold-savings/api/utils"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/config"
	"net/http"

	"github.com/google/uuid"
)

type Admin struct {
	queries *db.Queries
	config  *config.Config
}

func NewAdminService(q *db.Queries, c *config.Config) *Admin {
	return &Admin{q, c}
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

func (a *Admin) ApprovePayment(ctx context.Context, id int32, status, reason string) error {
	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "deposit",
		Reason: sql.NullString{String: reason, Valid: reason != ""},
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update transaction status: " + err.Error())
	}

	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	totalSavings := user.TotalSavings + trans.Amount
	params := db.UpdateUserTotalSavingsParams{
		ID:           user.ID,
		TotalSavings: totalSavings,
	}
	if err := a.queries.UpdateUserTotalSavings(ctx, params); err != nil {
		return errors.New("failed to update user total savings: " + err.Error())
	}

	// Optionally, you can add logic to notify the user about the status update
	emailbody, _ := utils.RenderEmailTemplate("templates/transaction_update.html", map[string]any{
		"Username":      user.Username,
		"Status":        status,
		"Amount":        trans.Amount,
		"Reason":        reason,
		"TransactionID": trans.ID,
		"Type":          trans.Type,
	})
	plunk := utils.Plunk{
		HttpClient: http.DefaultClient,
		Config:     a.config,
	}
	err = plunk.SendEmail(user.Email, "Transaction Status Update", emailbody)
	if err != nil {
		return errors.New("failed to send email notification: " + err.Error())
	}
	return nil
}

func (a *Admin) ApproveInvestment(ctx context.Context, id uuid.UUID, status string) error {
	args := db.UpdateInvestmentStatusParams{
		ID:     id,
		Status: status,
	}
	if err := a.queries.UpdateInvestmentStatus(ctx, args); err != nil {
		return errors.New("failed to update investment status: " + err.Error())
	}

	inv, err := a.queries.GetInvestmentByID(ctx, id)
	if err != nil {
		return errors.New("failed to get investment: " + err.Error())
	}

	userID, err := a.queries.GetUserFromInestmentID(ctx, id)
	if err != nil {
		return errors.New("failed to get user from investment ID: " + err.Error())
	}

	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	arg := db.UpdateUserTotalInvestmentBalanceParams{
		ID:                    user.ID,
		TotalInvestmentAmount: user.TotalInvestmentAmount + inv.Amount,
	}
	if err := a.queries.UpdateUserTotalInvestmentBalance(ctx, arg); err != nil {
		return errors.New("failed to update user total investment balance: " + err.Error())
	}

	return nil
}
