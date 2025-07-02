-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;


-- name: CreateFamily :exec
INSERT INTO families (name, description) VALUES (?, ?);

-- name: GetLastFiveFamilies :many
SELECT * FROM families ORDER BY created_at DESC LIMIT 5;