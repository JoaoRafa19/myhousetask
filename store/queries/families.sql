-- name: GetFamiliesByUserID :many
SELECT f.* 
FROM families f
JOIN family_members fm ON f.id = fm.family_id
WHERE fm.user_id = ?;
