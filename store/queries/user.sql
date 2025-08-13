-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByID :one
select * from users where id = ?;

-- name: CreateUser :exec
INSERT INTO users (id, name, email, password_hash)
VALUES (?, ?, ?, ?);
