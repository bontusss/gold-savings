-- name: CreateAdmin :one
INSERT INTO admins (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, created_at;

-- name: GetAdminByEmail :one
SELECT * FROM admins
WHERE email = $1 LIMIT 1;

-- name: CreateInvestmentPlan :one
INSERT INTO investment_plans (name, interest_rate, min_amount, max_amount)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAllInvestmentPlans :many
SELECT * FROM investment_plans;

-- name: GetInvestmentPlanByID :one
SELECT * FROM investment_plans WHERE id = $1;

-- name: DeleteInvestmentPlan :exec
DELETE FROM investment_plans
WHERE id = $1;

