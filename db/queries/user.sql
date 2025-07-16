
-- name: CreateUser :one
INSERT INTO users (email, password_hash, username, phone, reference_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUserStatus :exec
UPDATE users SET is_active = $2, updated_at = NOW() WHERE id = $1;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: SearchUsers :many
SELECT * FROM users
WHERE
    (first_name ILIKE '%' || $1 || '%' OR
     last_name ILIKE '%' || $1 || '%' OR
     email ILIKE '%' || $1 || '%' OR
     phone ILIKE '%' || $1 || '%')
ORDER BY created_at DESC;

-- name: CountActiveUsers :one
SELECT COUNT(*)
FROM users
WHERE is_active = TRUE;

-- name: SetUserEmailVerification :exec
UPDATE users
SET verification_code = $2,
    verification_expires_at = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: MarkUserEmailVerified :exec
UPDATE users
SET email_verified = $2,
    verification_code = $3,
    verification_expires_at = NULL,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserTotalSavings :exec
UPDATE users
SET total_savings = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetUserTotalSavings :one
SELECT total_savings
FROM users
WHERE id = $1;

-- name: UpdateUserTotalInvestmentBalance :exec
UPDATE users
SET total_investment_amount = $2,
    updated_at = NOW()
WHERE id = $1;