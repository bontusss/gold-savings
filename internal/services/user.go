package services

import (
	"context"
	"database/sql"
	"fmt"
	db "gold-savings/db/sqlc"

	"gold-savings/internal/auth"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserService struct {
	queries *db.Queries
}

// ErrInvestmentAmountOutOfRange is returned when the investment amount is not within the allowed range.
var ErrInvestmentAmountOutOfRange = fmt.Errorf("investment amount is out of allowed range")

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

type PayoutRequestType string
type PayoutRequestCategory string
type TransactionStatus string
type InvestmentStatus string

const (
	depositPayout        PayoutRequestCategory = "deposit"
	SavingsRequest       PayoutRequestType     = "savings"
	InvestmentRequest    PayoutRequestType     = "investment"
	PendingTransaction   TransactionStatus     = "pending"
	DeclinedTransaction  TransactionStatus     = "declined"
	WithdrawnTransaction TransactionStatus     = "withdrawal"
	DepositTransaction   TransactionStatus     = "deposit"
	PendingInvestment    InvestmentStatus      = "pending"
	ActiveInvestment     InvestmentStatus      = "active"
	CompletedInvestment  InvestmentStatus      = "completed"
)

func (s *UserService) GetUserByEmail(email string) (*db.User, error) {
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

func (s *UserService) GetInvestmentPlans(ctx context.Context) (*[]db.InvestmentPlan, error) {
	plans, err := s.queries.GetAllInvestmentPlans(ctx)
	if err != nil {
		return nil, err
	}
	return &plans, nil
}

func (s *UserService) CreateSavingsPaymentRequest(ctx context.Context, userID uuid.UUID, amount decimal.Decimal, bankName, account_name, account_number string) (*db.PayoutRequest, error) {
	args := db.CreatePayoutRequestParams{
		UserID:      userID,
		Amount:      amount.String(),
		BankName:    bankName,
		AccountName: account_name,
		Type:        string(SavingsRequest),
		Category:    string(DepositTransaction),
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *UserService) CreateInvestmentPaymentRequest(ctx context.Context, userID uuid.UUID, investment_id uuid.NullUUID, amount decimal.Decimal, bankName, account_name string) (*db.PayoutRequest, error) {
	args := db.CreatePayoutRequestParams{
		UserID:      userID,
		BankName:    bankName,
		AccountName: account_name,
		InvestmentID: investment_id,
		Type:        string(InvestmentRequest),
		Category:    string(DepositTransaction),
		Amount:      amount.String(),
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *UserService) CreateTransaction(ctx context.Context, userID uuid.UUID, amount decimal.Decimal, Ttype string) (*db.Transaction, error) {
	args := db.CreateTransactionParams{
		UserID: userID,
		Amount: amount.String(),
		Type:   Ttype,
		Status: string(PendingTransaction),
		Reason: sql.NullString{String: "", Valid: false},
	}
	transaction, err := s.queries.CreateTransaction(ctx, args)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *UserService) GetUserSavingsTransactions(ctx context.Context, userID uuid.UUID) (*[]db.Transaction, error) {
	transactions, err := s.queries.ListUserSavingsTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (s *UserService) GetUserInvestmentTransactions(ctx context.Context, userID uuid.UUID) (*[]db.Transaction, error) {
	transactions, err := s.queries.ListUserInvestmentTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (s *UserService) CreateInvestment(ctx context.Context, userID uuid.UUID, planID int32, amount decimal.Decimal) (*db.Investment, error) {
	iPlan, err := s.queries.GetInvestmentPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	IRefID := auth.GenerateReferenceID(iPlan.Name) // Generate a unique reference ID for the investment

	args := db.CreateInvestmentParams{
		UserID:       userID,
		PlanID:       planID,
		Amount:       amount.String(),
		Status:       string(PendingInvestment),
		ReferenceID:  IRefID,
		InterestRate: iPlan.InterestRate,
		Interest:     decimal.NewFromFloat(0).String(), // Initial interest is 0, will be calculated later
	}

	minAmount, err := decimal.NewFromString(iPlan.MinAmount)
	if err != nil {
		return nil, err
	}
	maxAmount, err := decimal.NewFromString(iPlan.MaxAmount)
	if err != nil {
		return nil, err
	}

	if amount.LessThan(minAmount) || amount.GreaterThan(maxAmount) {
		return nil, ErrInvestmentAmountOutOfRange
	}

	investment, err := s.queries.CreateInvestment(ctx, args)
	if err != nil {
		return nil, err
	}
	return &investment, nil
}

func (s *UserService) GetUserInvestments(ctx context.Context, userID uuid.UUID) (*[]db.Investment, error) {
	investments, err := s.queries.ListInvestmentsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &investments, nil
}

func (s *UserService) GetUserInvestmentsWithPlan(ctx context.Context, userID uuid.UUID) ([]db.ListUserInvestmentsWithPlanRow, error) {
	investments, err := s.queries.ListUserInvestmentsWithPlan(ctx, userID)
	if err != nil {
		return nil, err
	}
	return investments, nil
}

func (s *UserService) GetInvestmentPlanByID(ctx context.Context, planID int32) (*db.InvestmentPlan, error) {
	plan, err := s.queries.GetInvestmentPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *UserService) GetInvestmentByRefCode(ctx context.Context, refCode string) (*db.Investment, error) {
	pr, err := s.queries.GetInvestmentByRefCode(ctx, refCode)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}
