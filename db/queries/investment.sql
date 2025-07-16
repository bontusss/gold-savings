-- name: CreateInvestment :one
INSERT INTO investments (
  user_id, plan_id, amount, interest, interest_rate, status, reference_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetInvestmentByID :one
SELECT * FROM investments
WHERE id = $1;

-- name: GetUserFromInestmentID :one
SELECT user_id FROM investments
WHERE id = $1;

-- name: GetInvestmentByRefCode :one
SELECT * FROM investments
WHERE reference_id = $1;

-- name: ListUserInvestmentsWithPlan :many
SELECT
    i.*,
    p.name AS plan_name
FROM investments i
JOIN investment_plans p ON i.plan_id = p.id
WHERE i.user_id = $1
ORDER BY i.created_at DESC;

-- name: ListInvestmentsByUser :many
SELECT * FROM investments
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateInvestmentInterest :exec
UPDATE investments
SET interest = $2, updated_at = NOW()
WHERE id = $1;

-- name: ListInvestmentsByPlan :many
SELECT * FROM investments
WHERE plan_id = $1
ORDER BY created_at DESC;

-- name: UpdateInvestmentStatus :exec
UPDATE investments
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteInvestment :exec
DELETE FROM investments
WHERE id = $1;

-- name: CreateSavings :one
INSERT INTO savings (
  user_id, amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetSavingsByID :one
SELECT * FROM savings
WHERE id = $1;

-- name: ListSavingsByUserID :many
SELECT * FROM savings
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateSavingsAmount :one
UPDATE savings
SET amount = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteSavings :exec
DELETE FROM savings
WHERE id = $1;