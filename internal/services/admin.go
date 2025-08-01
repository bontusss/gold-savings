package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gold-savings/api/utils"
	db "gold-savings/db/sqlc"
	"gold-savings/internal/config"
	"log"
	"net/http"
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
	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "approved",
		Reason: sql.NullString{String: reason, Valid: reason != ""},
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update transaction status [savings payment request]: " + err.Error())
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	// Update user savings
	totalSavings := user.TotalSavings + trans.Amount
	params := db.UpdateUserTotalSavingsParams{
		ID:           user.ID,
		TotalSavings: totalSavings,
	}
	if err := a.queries.UpdateUserTotalSavings(ctx, params); err != nil {
		return errors.New("failed to update user total savings: " + err.Error())
	}

	// Update user tokens
	tokenParams := db.UpdateUserTokensParams{
		ID:          userID,
		TotalTokens: user.TotalTokens + int32(a.config.TOKENPERSAVINGS),
	}
	if err = a.queries.UpdateUserTokens(ctx, tokenParams); err != nil {
		fmt.Printf("failed to update user token: " + err.Error())
	}

	// Send confirmation email
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
		fmt.Printf("failed to send email notification: " + err.Error())
	}
	return nil
}

func (a *Admin) ApproveWithdrawal(ctx context.Context, id int32, status, reason string) error {
	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "approved",
		Reason: sql.NullString{String: reason, Valid: reason != ""},
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update transaction status: " + err.Error())
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	if trans.Amount > user.TotalSavings {
		return errors.New("Insufficient funds")
	}

	// Update user savings
	totalSavings := user.TotalSavings - trans.Amount
	params := db.UpdateUserTotalSavingsParams{
		ID:           user.ID,
		TotalSavings: totalSavings,
	}
	if err := a.queries.UpdateUserTotalSavings(ctx, params); err != nil {
		return errors.New("failed to update user total savings: " + err.Error())
	}

	// Update user tokens
	tokenParams := db.UpdateUserTokensParams{
		ID:          userID,
		TotalTokens: user.TotalTokens + int32(a.config.TOKENPERSAVINGS),
	}
	if err = a.queries.UpdateUserTokens(ctx, tokenParams); err != nil {
		return errors.New("failed to update user token: " + err.Error())
	}

	// Send confirmation email
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
		fmt.Printf("failed to send email notification: " + err.Error())
	}
	return nil
}

func (a *Admin) ApproveInvestment(ctx context.Context, id int32, status, reason string) error {
	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "approved",
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update Transaction status: " + err.Error())
	}

	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	investment, err := a.queries.GetInvestmentByID(ctx, trans.InvestmentID.Int32)
	if err != nil {
		return fmt.Errorf("error getting investment id: %v", err)
	}
	iParam := db.UpdateInvestmentStatusParams{
		ID:     investment.ID,
		Status: "active",
	}

	err = a.queries.UpdateInvestmentStatus(ctx, iParam)
	if err != nil {
		return fmt.Errorf("error updating investment status: %v", err)
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	arg := db.UpdateUserTotalInvestmentBalanceParams{
		ID:                    user.ID,
		TotalInvestmentAmount: user.TotalInvestmentAmount + trans.Amount,
	}
	if err := a.queries.UpdateUserTotalInvestmentBalance(ctx, arg); err != nil {
		return errors.New("failed to update user total investment balance: " + err.Error())
	}

	// Update user tokens
	tokenParams := db.UpdateUserTokensParams{
		ID:          userID,
		TotalTokens: user.TotalTokens + int32(a.config.TOKENPERINVESTMENT),
	}
	if err = a.queries.UpdateUserTokens(ctx, tokenParams); err != nil {
		return errors.New("failed to update user token: " + err.Error())
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

func (a *Admin) ApproveInvestmentWithdraw(ctx context.Context, id int32, status, reason string) error {
	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "approved",
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update Transaction status: " + err.Error())
	}

	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	investment, err := a.queries.GetInvestmentByID(ctx, trans.InvestmentID.Int32)
	if err != nil {
		return fmt.Errorf("error getting investment id: %v", err)
	}
	iParam := db.UpdateInvestmentStatusParams{
		ID:     investment.ID,
		Status: "active",
	}

	err = a.queries.UpdateInvestmentStatus(ctx, iParam)
	if err != nil {
		return fmt.Errorf("error updating investment status: %v", err)
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	arg := db.UpdateUserTotalInvestmentBalanceParams{
		ID:                    user.ID,
		TotalInvestmentAmount: user.TotalInvestmentAmount - trans.Amount,
	}
	if err := a.queries.UpdateUserTotalInvestmentBalance(ctx, arg); err != nil {
		return errors.New("failed to update user total investment balance: " + err.Error())
	}

	// // Update user tokens
	// tokenParams := db.UpdateUserTokensParams{
	// 	ID: userID,
	// 	TotalTokens: user.TotalTokens + int32(a.config.TOKENPERINVESTMENT),
	// }
	// if err = a.queries.UpdateUserTokens(ctx,tokenParams); err != nil {
	// 	return errors.New("failed to update user token: " + err.Error())
	// }

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
		fmt.Printf("failed to send email notification: " + err.Error())
	}

	return nil
}

func (a *Admin) ApproveTokenRequest(ctx context.Context, id int32, status, reason string) error {
	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "approved",
		Reason: sql.NullString{String: reason, Valid: reason != ""},
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update transaction status [savings payment request]: " + err.Error())
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	// Update user savings
	totalTokens := user.TotalTokens - trans.Amount
	params := db.UpdateUserTokensParams{
		ID:          user.ID,
		TotalTokens: totalTokens,
	}
	if err := a.queries.UpdateUserTokens(ctx, params); err != nil {
		return errors.New("failed to update user total savings: " + err.Error())
	}

	// Send confirmation email
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
		fmt.Printf("failed to send email notification: " + err.Error())
	}
	return nil
}

func (a *Admin) DeclinePayment(ctx context.Context, id int32, status, reason string) error {
	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "declined",
		Reason: sql.NullString{String: reason, Valid: reason != ""},
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update transaction status: " + err.Error())
	}

	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}
	log.Printf("transaction status is %s", trans.Status)

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	// Optionally, you can add logic to notify the user about the status update
	emailbody, _ := utils.RenderEmailTemplate("templates/transaction_update.html", map[string]any{
		"Username":      user.Username,
		"Status":        trans.Status,
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
		fmt.Printf("failed to send email notification: " + err.Error())
	}
	return nil
}

func (a *Admin) DeclineInvestment(ctx context.Context, id int32, status, reason string) error {
	args := db.UpdateTransactionStatusParams{
		ID:     id,
		Status: "declined",
	}
	if err := a.queries.UpdateTransactionStatus(ctx, args); err != nil {
		return errors.New("failed to update Transaction status: " + err.Error())
	}

	trans, err := a.queries.GetTransactionByID(ctx, id)
	if err != nil {
		return errors.New("failed to get transaction: " + err.Error())
	}

	investment, err := a.queries.GetInvestmentByID(ctx, trans.InvestmentID.Int32)
	if err != nil {
		return fmt.Errorf("error getting investment id: %v", err)
	}
	iParam := db.UpdateInvestmentStatusParams{
		ID:     investment.ID,
		Status: "cancelled",
	}

	err = a.queries.UpdateInvestmentStatus(ctx, iParam)
	if err != nil {
		return fmt.Errorf("error updating investment status: %v", err)
	}

	userID, err := a.queries.GetUserFromTransactionID(ctx, id)
	if err != nil {
		return err
	}
	user, err := a.queries.GetUser(ctx, userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	// Optionally, you can add logic to notify the user about the status update
	emailbody, _ := utils.RenderEmailTemplate("templates/transaction_update.html", map[string]any{
		"Username":      user.Username,
		"Status":        trans.Status,
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
		fmt.Printf("failed to send email notification: " + err.Error())
	}

	return nil
}
