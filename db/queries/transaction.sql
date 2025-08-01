-- name: CreateTransaction :one
INSERT INTO transactions (
  user_id, amount, type, investment_id, status, reason, category
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = $1;

-- name: GetUserFromTransactionID :one
SELECT user_id FROM transactions
WHERE id = $1;

-- name: ListTransactionsByUserID :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListUserSavingsTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
  AND category = 'savings'
ORDER BY created_at DESC;

-- name: ListUserInvestmentTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
  AND type = 'investment'
ORDER BY created_at DESC;

-- name: ListTransactionsByCategory :many
SELECT
  transactions.*,
  users.username
FROM transactions
JOIN users ON users.id = transactions.user_id
WHERE transactions.category = $1
ORDER BY transactions.created_at DESC;


-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET status = $2,
    reason = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetRejectedTransactions :many
SELECT * FROM transactions
WHERE status = 'rejected';

-- name: ListPendingDepositTransactionsWithUser :many
SELECT
  t.id,
  t.user_id,
  t.amount,
  t.type,
  t.status,
  t.reason,
  t.created_at,
  t.updated_at,
  u.username,
  u.email
FROM transactions t
JOIN users u ON t.user_id = u.id
WHERE t.status = 'pending'
AND t.type = 'deposit'
ORDER BY t.created_at DESC;

-- name: ListPendingWithdrawalTransactionsWithUser :many
SELECT
  t.id,
  t.user_id,
  t.amount,
  t.type,
  t.status,
  t.reason,
  t.created_at,
  t.updated_at,
  u.username,
  u.email
FROM transactions t
JOIN users u ON t.user_id = u.id
WHERE t.status = 'pending'
AND t.type = 'withdrawal'
ORDER BY t.created_at DESC;

-- name: CreatePayoutRequest :one
INSERT INTO payout_requests (
  user_id, account_name, bank_name,account_number, investment_id, type, category, amount, phone_number
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetPayoutRequestByID :one
SELECT * FROM payout_requests
WHERE id = $1;

-- name: ListPayoutRequestsByUserID :many
SELECT * FROM payout_requests
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListPayoutRequestsByCategory :many
SELECT * FROM payout_requests
WHERE category = $1
ORDER BY created_at DESC;

-- name: ListAllPayoutRequests :many
SELECT * FROM payout_requests;

-- name: ListPayoutRequestsByType :many
SELECT * FROM payout_requests
WHERE type = $1
ORDER BY created_at DESC;

-- name: UpdatePayoutRequestAmount :exec
UPDATE payout_requests
SET amount = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

