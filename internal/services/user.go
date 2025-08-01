package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gold-savings/api/utils"
	db "gold-savings/db/sqlc"
	"net/http"

	"gold-savings/internal/auth"
	"gold-savings/internal/config"

	"github.com/shopspring/decimal"
)

type UserService struct {
	queries *db.Queries
	config  *config.Config
}

// ErrInvestmentAmountOutOfRange is returned when the investment amount is not within the allowed range.
var ErrInvestmentAmountOutOfRange = fmt.Errorf("investment amount is out of allowed range")

func NewUserService(queries *db.Queries, c *config.Config) *UserService {
	return &UserService{
		queries: queries,
		config:  c,
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

func (s *UserService) UpdateStatus(userid int32, isActive bool) error {
	err := s.queries.UpdateUserStatus(context.Background(), db.UpdateUserStatusParams{
		ID:       userid,
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

func (s *UserService) CreateSavingsPaymentRequest(ctx context.Context, userID int32, amount decimal.Decimal, bankName, account_name string) (*db.PayoutRequest, error) {
	args := db.CreatePayoutRequestParams{
		UserID:      userID,
		Amount:      amount.String(),
		BankName:    sql.NullString{Valid: true, String: bankName},
		AccountName: sql.NullString{Valid: true, String: account_name},
		Type:        string(SavingsRequest),
		Category:    string(DepositTransaction),
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	plunk := utils.Plunk{
		HttpClient: http.DefaultClient,
		Config:     s.config,
	}
	emailbody, _ := utils.RenderEmailTemplate("templates/savings_payment_request.html", map[string]any{
		"Username":    user.Username,
		"Amount":      amount.String(),
		"BankName":    bankName,
		"AccountName": account_name,
		"Type":        payment.Type,
		"Category":    payment.Category,
		"Date":        payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		"ApproveURL":  fmt.Sprintf("%s/admin/api/approve-payment/%d", s.config.BaseURL, payment.ID),
		"RejectURL":   fmt.Sprintf("%s/admin/api/decline-payment/%d", s.config.BaseURL, payment.ID),
	})

	fmt.Printf("Processing payment ID: %d for user: %s\n", payment.ID, user.Username)

	emailbody2, _ := utils.RenderEmailTemplate("templates/user_payment_request.html", map[string]any{
		"Username":    user.Username,
		"Amount":      amount.String(),
		"BankName":    bankName,
		"AccountName": account_name,
		"Category":    payment.Category,
		"Type":        payment.Type,
		"Date":        payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	})
	fmt.Println("sending admin email ...")
	err = plunk.SendEmail(s.config.AdminEmail, "New Savings Payment Request", emailbody)
	if err != nil {
		fmt.Printf("failed to send admin email notification [savings payment request]: %v", err)
	}

	fmt.Println("sending user email ...")
	err = plunk.SendEmail(user.Email, "New Savings Payment Request [savings payment request]", emailbody2)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %v", err)
	}
	return &payment, nil
}

func (s *UserService) CreateSavingsWithdrawRequest(ctx context.Context, userID int32, amount decimal.Decimal, bankName, account_name, accountNumber string) (*db.PayoutRequest, error) {
	args := db.CreatePayoutRequestParams{
		UserID:        userID,
		Amount:        amount.String(),
		BankName:      sql.NullString{Valid: true, String: bankName},
		AccountName:   sql.NullString{Valid: true, String: account_name},
		AccountNumber: sql.NullString{Valid: true, String: accountNumber},
		Type:          string(SavingsRequest),
		Category:      string(WithdrawnTransaction),
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	plunk := utils.Plunk{
		HttpClient: http.DefaultClient,
		Config:     s.config,
	}
	emailbody, _ := utils.RenderEmailTemplate("templates/savings_payment_request.html", map[string]any{
		"Username":      user.Username,
		"Amount":        amount.String(),
		"BankName":      bankName,
		"AccountName":   account_name,
		"AccountNumber": accountNumber,
		"Type":          payment.Type,
		"Category":      payment.Category,
		"Date":          payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		"ApproveURL":    fmt.Sprintf("%s/admin/api/approve-payment/%d", s.config.BaseURL, payment.ID),
		"RejectURL":     fmt.Sprintf("%s/admin/api/decline-payment/%d", s.config.BaseURL, payment.ID),
	})

	fmt.Printf("Processing payment ID: %d for user: %s\n", payment.ID, user.Username)

	emailbody2, _ := utils.RenderEmailTemplate("templates/user_payment_request.html", map[string]any{
		"Username":      user.Username,
		"Amount":        amount.String(),
		"BankName":      bankName,
		"AccountName":   account_name,
		"AccountNumber": accountNumber,
		"Type":          payment.Type,
		"Category":      payment.Category,
		"Date":          payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	})
	fmt.Println("sending admin email ...")
	err = plunk.SendEmail(s.config.AdminEmail, "New Savings Payment Request", emailbody)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %w", err)
	}

	fmt.Println("sending user email ...")
	err = plunk.SendEmail(user.Email, "New Savings Payment Request", emailbody2)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %w", err)
	}
	return &payment, nil
}

func (s *UserService) CreateInvestmentPaymentRequest(ctx context.Context, userID int32, investment_id int32, amount decimal.Decimal, bankName, account_name string) (*db.PayoutRequest, error) {
	args := db.CreatePayoutRequestParams{
		UserID:       userID,
		BankName:     sql.NullString{Valid: true, String: bankName},
		AccountName:  sql.NullString{Valid: true, String: account_name},
		InvestmentID: sql.NullInt32{Valid: true, Int32: investment_id},
		Type:         string(InvestmentRequest),
		Category:     string(DepositTransaction),
		Amount:       amount.String(),
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	plunk := utils.Plunk{
		HttpClient: http.DefaultClient,
		Config:     s.config,
	}
	emailbody, _ := utils.RenderEmailTemplate("templates/savings_payment_request.html", map[string]any{
		"Username":      user.Username,
		"Amount":        amount.String(),
		"BankName":      bankName,
		"AccountName":   account_name,
		"AccountNumber": "-",
		"Type":          payment.Type,
		"Date":          payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		"ApproveURL":    fmt.Sprintf("%s/admin/api/approve-investment/%d", s.config.BaseURL, payment.ID),
		"RejectURL":     fmt.Sprintf("%s/admin/api/decline-investment/%d", s.config.BaseURL, payment.ID),
	})

	fmt.Printf("Processing payment ID: %d for user: %s\n", payment.ID, user.Username)

	emailbody2, _ := utils.RenderEmailTemplate("templates/user_payment_request.html", map[string]any{
		"Username":    user.Username,
		"Amount":      amount.String(),
		"BankName":    bankName,
		"AccountName": account_name,
		"Type":        payment.Type,
		"Date":        payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	})
	err = plunk.SendEmail(s.config.AdminEmail, "New Investment Payment Request", emailbody)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %v", err)
	}
	err = plunk.SendEmail(user.Email, "New Investment Payment Request", emailbody2)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %v", err)
	}
	return &payment, nil
}

func (s *UserService) CreateInvestmentWithdrawRequest(ctx context.Context, userID int32, investment_id int32, amount decimal.Decimal, bankName, account_name, account_number string) (*db.PayoutRequest, error) {
	args := db.CreatePayoutRequestParams{
		UserID:        userID,
		BankName:      sql.NullString{Valid: true, String: bankName},
		AccountName:   sql.NullString{Valid: true, String: account_name},
		AccountNumber: sql.NullString{Valid: true, String: account_number},
		InvestmentID:  sql.NullInt32{Valid: true, Int32: investment_id},
		Type:          string(InvestmentRequest),
		Category:      string(WithdrawnTransaction),
		Amount:        amount.String(),
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	plunk := utils.Plunk{
		HttpClient: http.DefaultClient,
		Config:     s.config,
	}
	emailbody, _ := utils.RenderEmailTemplate("templates/savings_payment_request.html", map[string]any{
		"Username":      user.Username,
		"Amount":        amount.String(),
		"BankName":      bankName,
		"AccountName":   account_name,
		"AccountNumber": account_number,
		"Category":      payment.Category,
		"Type":          payment.Type,
		"Date":          payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		"ApproveURL":    fmt.Sprintf("%s/admin/api/approve-investment/%d", s.config.BaseURL, payment.ID),
		"RejectURL":     fmt.Sprintf("%s/admin/api/decline-investment/%d", s.config.BaseURL, payment.ID),
	})

	fmt.Printf("Processing payment ID: %d for user: %s\n", payment.ID, user.Username)

	emailbody2, _ := utils.RenderEmailTemplate("templates/user_payment_request.html", map[string]any{
		"Username":      user.Username,
		"Amount":        amount.String(),
		"BankName":      bankName,
		"AccountName":   account_name,
		"AccountNumber": account_number,
		"Category":      payment.Category,
		"Type":          payment.Type,
		"Date":          payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	})
	err = plunk.SendEmail(s.config.AdminEmail, "New Investment Payment Request", emailbody)
	if err != nil {
		return nil, fmt.Errorf("failed to send admin email notification: %w", err)
	}
	err = plunk.SendEmail(user.Email, "New Investment Payment Request", emailbody2)
	if err != nil {
		return nil, fmt.Errorf("failed to send admin email notification: %w", err)
	}
	return &payment, nil
}

func (s *UserService) CreateTransaction(ctx context.Context, userID int32, investmentID sql.NullInt32, amount decimal.Decimal, Ttype, category string) (*db.Transaction, error) {
	args := db.CreateTransactionParams{
		UserID:       userID,
		Amount:       int32(amount.IntPart()),
		Type:         Ttype,
		InvestmentID: investmentID,
		Category:     category,
		Status:       string(PendingTransaction),
		Reason:       sql.NullString{String: "", Valid: false},
	}
	transaction, err := s.queries.CreateTransaction(ctx, args)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *UserService) GetUserSavingsTransactions(ctx context.Context, userID int32) (*[]db.Transaction, error) {
	transactions, err := s.queries.ListUserSavingsTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (s *UserService) GetUserInvestmentTransactions(ctx context.Context, userID int32) (*[]db.Transaction, error) {
	transactions, err := s.queries.ListUserInvestmentTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (s *UserService) CreateInvestment(ctx context.Context, userID int32, planID int32, amount decimal.Decimal) (*db.Investment, error) {
	iPlan, err := s.queries.GetInvestmentPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	IRefID := auth.GenerateReferenceID(iPlan.Name) // Generate a unique reference ID for the investment

	args := db.CreateInvestmentParams{
		UserID:       userID,
		PlanID:       planID,
		Amount:       int32(amount.IntPart()),
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

func (s *UserService) GetUserInvestments(ctx context.Context, userID int32) (*[]db.Investment, error) {
	investments, err := s.queries.ListInvestmentsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &investments, nil
}

func (s *UserService) GetUserInvestmentsWithPlan(ctx context.Context, userID int32) ([]db.ListUserInvestmentsWithPlanRow, error) {
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

func (s *UserService) GetUserByID(ctx context.Context, userID int32) (*db.User, error) {
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserSavingsBalance(ctx context.Context, userID int32) (int32, error) {
	totalSavings, err := s.queries.GetUserTotalSavings(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("error getting user savings balance: %v", err)
	}
	return totalSavings, nil
}

func (s *UserService) GetUserInvestmentBalance(ctx context.Context, userID int32) (int32, error) {
	totalSavings, err := s.queries.GetUserInvestmentBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("error getting user savings balance: %v", err)
	}
	return totalSavings, nil
}

func (s *UserService) UpdateEmailAndUsername(ctx context.Context, email, username string, id int32) error {
	args := db.UpdateUsernameEmailParams{
		ID:       id,
		Username: username,
		Email:    email,
	}
	if err := s.queries.UpdateUsernameEmail(ctx, args); err != nil {
		return fmt.Errorf("error updating username and email: %v", err)
	}
	return nil
}

func (s *UserService) UpdateEmailOrUsername(ctx context.Context, email, username *string, id int32) error {
	args := db.UpdateUsernameEmailPartialParams{
		ID: id,
		Email: sql.NullString{
			String: deref(email),
			Valid:  email != nil,
		},
		Username: sql.NullString{
			String: deref(username),
			Valid:  username != nil,
		},
	}

	if err := s.queries.UpdateUsernameEmailPartial(ctx, args); err != nil {
		return fmt.Errorf("error updating username and email: %v", err)
	}
	return nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *UserService) ListUserReferrals(ctx context.Context, userID int32) (*[]db.ListReferralsByInviterRow, error) {
	refs, err := s.queries.ListReferralsByInviter(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &refs, nil
}

func (s *UserService) CreateTokenRedeemRequest(ctx context.Context, userID int32, amount decimal.Decimal, phoneNumber string) (*db.PayoutRequest, error) {
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if amount.LessThan(decimal.NewFromInt(int64(s.config.TokenThreshold))) {
		return nil, fmt.Errorf("You can redeem token from %d and above", s.config.TokenThreshold)
	}

	if amount.GreaterThan(decimal.NewFromInt(int64(user.TotalTokens))) {
		return nil, errors.New("insufficient token balance")
	}

	args := db.CreatePayoutRequestParams{
		UserID:      userID,
		Amount:      amount.String(),
		PhoneNumber: sql.NullString{Valid: true, String: phoneNumber},
		Type:        "token",
		Category:    "token",
	}
	payment, err := s.queries.CreatePayoutRequest(ctx, args)
	if err != nil {
		return nil, err
	}

	plunk := utils.Plunk{
		HttpClient: http.DefaultClient,
		Config:     s.config,
	}
	emailbody, _ := utils.RenderEmailTemplate("templates/savings_payment_request.html", map[string]any{
		"Username": user.Username,
		"Amount":   amount.String(),
		"Phone":    phoneNumber,
		// "BankName":    bankName,
		// "AccountName": account_name,
		"Type":       payment.Type,
		"Category":   payment.Category,
		"Date":       payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		"ApproveURL": fmt.Sprintf("%s/admin/api/approve-payment/%d", s.config.BaseURL, payment.ID),
		"RejectURL":  fmt.Sprintf("%s/admin/api/decline-payment/%d", s.config.BaseURL, payment.ID),
	})

	fmt.Printf("Processing payment ID: %d for user: %s\n", payment.ID, user.Username)

	emailbody2, _ := utils.RenderEmailTemplate("templates/user_payment_request.html", map[string]any{
		"Username": user.Username,
		"Amount":   amount.String(),
		// "BankName":    bankName,
		// "AccountName": account_name,
		"Category": payment.Category,
		"Type":     payment.Type,
		"Date":     payment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	})
	fmt.Println("sending admin email ...")
	err = plunk.SendEmail(s.config.AdminEmail, "New Savings Payment Request", emailbody)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %v", err)
	}

	fmt.Println("sending user email ...")
	err = plunk.SendEmail(user.Email, "New Savings Payment Request", emailbody2)
	if err != nil {
		fmt.Printf("failed to send admin email notification: %v", err)
	}
	return &payment, nil
}

func (s *UserService) GetUserTokens(ctx context.Context, userID int32) (int32, error) {
	tokens, err := s.queries.GetUserTokens(ctx, userID)
	if err != nil {
		return 0, err
	}
	return tokens, nil
}
