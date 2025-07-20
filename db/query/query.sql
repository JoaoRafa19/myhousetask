-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;


-- name: CreateFamily :execresult
INSERT INTO families (name, description, owner_id) VALUES (?, ?, ?);

-- name: ListFamiliesForUser :many
SELECT f.*
FROM families f
         JOIN family_members fm ON f.id = fm.family_id
WHERE fm.user_id = ?
ORDER BY f.created_at DESC;


-- name: DashboardPage :many
SELECT
    f.id as id_familia,
    f.name as nome_familia,
    DATE_FORMAT(f.created_at, '%Y-%m-%d %H:%i:%s') as created_at,
    f.is_active as status,
    (SELECT COUNT(*) FROM family_members WHERE family_id = f.id) as total_membros
FROM families f
WHERE f.id IN (
    SELECT family_id FROM family_members fm WHERE fm.user_id = ?
)
ORDER BY f.created_at DESC
LIMIT 5;

-- name: CountFamilies :one
SELECT count(*)
FROM family_members fm
WHERE fm.user_id = ?;

-- name: CountUsersFamilyMembers :one
SELECT COUNT(DISTINCT fm2.user_id)
FROM family_members AS fm1
         JOIN family_members AS fm2 ON fm1.family_id = fm2.family_id
WHERE fm1.user_id = ?;

-- name: GetUserByID :one
    select * from users where id = ?;

-- name: CountTasksCompletedToday :one
SELECT count(*)
FROM tasks t
         JOIN family_members fm ON t.family_id = fm.family_id
WHERE
    fm.user_id = ? AND
    t.status = 'completed' AND
    DATE(t.completed_at) = CURDATE();


-- name: CountTasksPending :one
SELECT count(*) FROM tasks WHERE status = 'pending';

-- name: ListRecentFamilies :many
SELECT
    f.id as id_familia,
    f.name as nome_familia,
    DATE_FORMAT(f.created_at, '%Y-%m-%d %H:%i:%s') as created_at,
    f.is_active as status,
    COUNT(fm.id) as total_membros
from families f
LEFT JOIN family_members fm on fm.family_id = f.id
LEFT JOIN users u on fm.user_id = u.id
GROUP BY f.id, f.name, f.created_at, f.is_active
ORDER BY f.created_at DESC
LIMIT 5;


-- name: GetWeeklyTaskCompletionStats :many
SELECT
    DATE(completed_at) as completion_date,
    COUNT(*) as completed_count
FROM tasks
WHERE
    status = 'completed' AND
    completed_at >= CURDATE() - INTERVAL 7 DAY
GROUP BY completion_date
ORDER BY completion_date DESC;


-- name: CreateUser :exec
INSERT INTO users (id, name, email, password_hash)
VALUES (?, ?, ?, ?);

-- name: CreateFamilyMember :exec
INSERT INTO family_members (id, family_id, user_id, role)
VALUES (?, ?, ?, ?);
