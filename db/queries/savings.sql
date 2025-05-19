-- queries.sql

-- name: CreateSavingsPlan :one
INSERT INTO savings_plans (
    user_id,
    plan_ref,
    target_amount,
    current_amount,
    duration_days,
    savings_frequency,
    savings_amount,
    status
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8
         ) RETURNING *;

-- name: GetSavingsPlanByID :one
SELECT * FROM savings_plans
WHERE id = $1;

-- name: GetSavingsPlanByPlanRef :one
SELECT * FROM savings_plans
WHERE plan_ref = $1;

-- name: ListSavingsPlansByUser :many
SELECT * FROM savings_plans
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateSavingsPlanAmount :one
UPDATE savings_plans
SET current_amount = $2,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: UpdateSavingsPlanStatus :one
UPDATE savings_plans
SET status = $2,
    updated_at = NOW()
WHERE id = $1
    RETURNING *;

-- name: ListActiveSavingsPlans :many
SELECT * FROM savings_plans
WHERE status = 'active'
  AND maturity_date > NOW()
ORDER BY maturity_date ASC;

-- name: DeleteSavingsPlan :exec
DELETE FROM savings_plans
WHERE id = $1;

-- name: ListAllSavingsPlans :many
SELECT * FROM savings_plans
ORDER BY created_at DESC;

-- name: SumActiveSavingsPlans :one
SELECT COALESCE(SUM(current_amount),0)::float8 AS total_active_savings
FROM savings_plans
WHERE status = 'active';