
-- name: CreateUser :one
INSERT INTO users (email, password_hash, username, phone, reference_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByReferenceID :one
SELECT * FROM users
WHERE reference_id = $1 LIMIT 1;

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

-- name: UpdateUsernameEmail :exec
UPDATE users
SET email = $2,
    username = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUsernameEmailPartial :exec
UPDATE users
SET
  email = COALESCE(sqlc.narg('email'), email),
  username = COALESCE(sqlc.narg('username'), username)
WHERE id = sqlc.arg('id');


-- name: UpdateUserTotalSavings :exec
UPDATE users
SET total_savings = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetUserTotalSavings :one
SELECT total_savings
FROM users
WHERE id = $1;

-- name: GetUserTotalTokeens :one
SELECT total_tokens
FROM users
WHERE id = $1;

-- name: GetUserTokens :one
SELECT total_tokens
FROM users
WHERE id = $1;

-- name: UpdateUserTokens :exec
UPDATE users
SET total_tokens = $2
WHERE id = $1;

-- name: GetUserInvestmentBalance :one
SELECT total_investment_amount
FROM users
WHERE id = $1;

-- name: UpdateUserTotalInvestmentBalance :exec
UPDATE users
SET total_investment_amount = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetTotalSavingsInApp :one
SELECT COALESCE(SUM(total_savings), 0)::INTEGER AS total FROM users;

-- name: GetTotalInvestmentInApp :one
SELECT COALESCE(SUM(total_investment_amount), 0)::INTEGER AS total FROM users;

-- name: GetTotalTokens :one
SELECT COALESCE(SUM(total_tokens), 0)::INTEGER AS total FROM users;

-- name: ListUsersByTotalSavingsDesc :many
SELECT *
FROM users
WHERE total_savings > 0
ORDER BY total_savings DESC
LIMIT 10;

-- name: CreateReferral :one
INSERT INTO referrals (inviter_id, invitee_id)
VALUES ($1, $2)
RETURNING id, inviter_id, invitee_id, created_at;

-- name: ListReferralsByInviter :many
SELECT r.id, r.inviter_id, r.invitee_id, r.created_at, u.email, u.username
FROM referrals r
JOIN users u ON u.id = r.invitee_id
WHERE r.inviter_id = $1
ORDER BY r.created_at DESC;
