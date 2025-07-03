-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;


-- name: CreateFamily :exec
INSERT INTO families (name, description) VALUES (?, ?);

-- name: GetLastFiveFamilies :many
SELECT * FROM families ORDER BY created_at DESC LIMIT 5;


-- name: DashboardPage :many 
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

-- name: CountFamilies :one
SELECT count(*) FROM families;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: CountTasksCompletedToday :one
SELECT count(*) FROM tasks WHERE status = 'completed' AND DATE(completed_at) = CURDATE();

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
