-- name: CreateUser :one
INSERT INTO user (
    id,
    username,
    email,
    email_confirm_code,
    password
) VALUES (
    ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM user WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM user WHERE username = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM user WHERE email = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE user
SET 
    username = COALESCE(?, username),
    email = COALESCE(?, email),
    password = COALESCE(?, password),
    email_confirmed_at = COALESCE(?, email_confirmed_at),
    email_confirm_code = COALESCE(?, email_confirm_code)
WHERE id = ?;

-- name: ConfirmEmail :exec
UPDATE user
SET 
    email_confirmed_at = CURRENT_TIMESTAMP,
    email_confirm_code = NULL
WHERE id = ? AND email_confirm_code = ?;

-- name: DeleteUser :exec
DELETE FROM user WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM user ORDER BY created_at DESC;
