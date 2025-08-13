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


-- name: CreateTask :exec
INSERT INTO tasks (
    id, family_id, title, description, status_id, created_by
) VALUES (
     ?, ?, ?, ?, ?, ?
);

-- name: GetTasksByFamily :many
SELECT
    t.*,
    ts.name as status_name
FROM tasks t
         JOIN task_status ts ON t.status_id = ts.id
WHERE t.family_id = ?;

-- name: GetTaskStatus :many
SELECT * FROM task_status;


-- name: GetTask :one
SELECT * from tasks where id = ?;