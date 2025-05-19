-- name: CreateAdminUser :one
INSERT INTO users (email, password_hash, is_admin, first_name, last_name, phone)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, is_admin, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, is_admin FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUserStatus :exec
UPDATE users SET is_active = $2, updated_at = NOW() WHERE id = $1;

-- name: SearchUsers :many
SELECT * FROM users 
WHERE 
    (first_name ILIKE '%' || $1 || '%' OR 
     last_name ILIKE '%' || $1 || '%' OR 
     email ILIKE '%' || $1 || '%' OR 
     phone ILIKE '%' || $1 || '%')
ORDER BY created_at DESC;

-- name: CountActiveNonAdminUsers :one
SELECT COUNT(*)
FROM users
WHERE is_active = TRUE AND is_admin = FALSE;