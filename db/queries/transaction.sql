-- name: CreateTransaction :one
INSERT INTO transactions (
  user_id, amount, type, status, reason
) VALUES (
  $1, $2, $3, $4, $5
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
  AND type = 'savings'
ORDER BY created_at DESC;

-- name: ListUserInvestmentTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
  AND type = 'investment'
ORDER BY created_at DESC;

-- name: ListTransactionsByType :many
SELECT * FROM transactions
WHERE type = $1
ORDER BY created_at DESC;

-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET status = $2,
    reason = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetRejectedTransactions :many
SELECT * FROM transactions
WHERE status = 'rejected';


-- name: CreatePayoutRequest :one
INSERT INTO payout_requests (
  user_id, account_name, bank_name, investment_id, type, category, amount
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
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

-- name: ListPayoutRequestsByType :many
SELECT * FROM payout_requests
WHERE type = $1
ORDER BY created_at DESC;

-- name: UpdatePayoutRequestAmount :exec
UPDATE payout_requests
SET amount = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

